// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// ofactest is a cli tool used for testing the Moov OFAC service.
//
// With no arguments the contaier runs tests against the production API.
// This tool requires an OAuth token provided by github.com/moov-io/api written
// to the local disk, but running apitest first will write this token.
//
// This tool can be used to query with custom searches:
//  $ go install ./cmd/ofactest
//  $ ofactest -local moh
//  2019/02/14 23:37:44.432334 main.go:44: Starting moov/ofactest v0.4.1-dev
//  2019/02/14 23:37:44.432366 main.go:60: [INFO] using http://localhost:8084 for address
//  2019/02/14 23:37:44.434534 main.go:76: [SUCCESS] ping
//  2019/02/14 23:37:44.435204 main.go:83: [SUCCESS] last download was: 3h45m58s ago
//  2019/02/14 23:37:44.440230 main.go:96: [SUCCESS] name search passed, query="moh"
//  2019/02/14 23:37:44.441506 main.go:104: [SUCCESS] added customer=24032 watch
//  2019/02/14 23:37:44.445473 main.go:118: [SUCCESS] alt name search passed
//  2019/02/14 23:37:44.449367 main.go:123: [SUCCESS] address search passed
//
// ofactest is not a stable tool. Please contact Moov developers if you intend to use this tool,
// otherwise we might change the tool (or remove it) without notice.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cardonator/ofac"
	moov "github.com/cardonator/ofac/client"
	"github.com/moov-io/base/http/bind"

	"github.com/antihax/optional"
)

var (
	defaultApiAddress = "https://api.moov.io"

	flagApiAddress = flag.String("address", defaultApiAddress, "Moov API address")
	flagLocal      = flag.Bool("local", false, "Use local HTTP addresses")
	flagWebhook    = flag.String("webhook", "https://moov.io/ofac", "Secure HTTP address for webhooks")
)

func main() {
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds | log.Lshortfile)
	log.Printf("Starting moov/ofactest %s", ofac.Version)

	conf := moov.NewConfiguration()
	conf.BasePath = getBasePath(*flagApiAddress, *flagLocal)

	conf.UserAgent = fmt.Sprintf("moov/ofactest:%s", ofac.Version)
	conf.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 1 * time.Minute,
		},
	}

	log.Printf("[INFO] using %s for address", conf.BasePath)

	// Read OAuth token and set on conf
	if v := os.Getenv("OAUTH_TOKEN"); v != "" {
		conf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", v))
	} else {
		if local := *flagLocal; !local {
			log.Fatal("[FAILURE] no OAuth token provided")
		}
	}

	// Setup OFAC API client
	api, ctx := moov.NewAPIClient(conf), context.TODO()

	// Ping OFAC
	if err := ping(ctx, api); err != nil {
		log.Fatal("[FAILURE] ping OFAC")
	} else {
		log.Println("[SUCCESS] ping")
	}

	// Check downloads
	if when, err := latestDownload(ctx, api); err != nil || when.IsZero() {
		log.Fatalf("[FAILURE] downloads: %v", err)
	} else {
		log.Printf("[SUCCESS] last download was: %v ago", time.Since(when).Truncate(1*time.Second))
	}

	query := "alh" // string that matches a lot of OFAC records
	if v := flag.Arg(0); v != "" {
		query = v
	}

	// Search queries
	sdn, err := searchByName(ctx, api, query)
	if err != nil {
		log.Fatalf("[FAILURE] problem searching SDNs: %v", err)
	} else {
		log.Printf("[SUCCESS] name search passed, query=%q", query)
	}

	// Add watch on the SDN
	if strings.EqualFold(sdn.SdnType, "individual") {
		if err := addCustomerWatch(ctx, api, sdn.EntityID, *flagWebhook); err != nil {
			log.Fatalf("[FAILURE] problem adding customer watch: %v", err)
		} else {
			log.Printf("[SUCCESS] added customer=%s watch", sdn.EntityID)
		}
	} else {
		if err := addCompanyWatch(ctx, api, sdn.EntityID, *flagWebhook); err != nil {
			log.Fatalf("[FAILURE] problem adding company watch: %v", err)
		} else {
			log.Printf("[SUCCESS] added company=%s watch", sdn.EntityID)
		}
	}

	// Load alt names and addresses
	if err := searchByAltName(ctx, api, query); err != nil {
		log.Fatalf("[FAILURE] problem searching Alt Names: %v", err)
	} else {
		log.Println("[SUCCESS] alt name search passed")
	}
	if err := searchByAddress(ctx, api, "St"); err != nil {
		log.Fatalf("[FAILURE] problem searching addresses: %v", err)
	} else {
		log.Println("[SUCCESS] address search passed")
	}
}

// getBasePath reads flagLocal and flagApiAddress to compute the HTTP address used for connecting with OFAC.
func getBasePath(address string, local bool) string {
	if local {
		// If '-local and -address <foo>' use <foo>
		if address != defaultApiAddress {
			return strings.TrimSuffix(address, "/")
		} else {
			return "http://localhost" + bind.HTTP("ofac")
		}
	} else {
		address = strings.TrimSuffix(address, "/")
		// -address isn't changed, so assume Moov's API (needs extra path added)
		if address == defaultApiAddress {
			return address + "/v1/ofac"
		}
		return address
	}
}

func ping(ctx context.Context, api *moov.APIClient) error {
	resp, err := api.OFACApi.Ping(ctx)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("ping error (stats code: %d): %v", resp.StatusCode, err)
	}
	return nil
}

func latestDownload(ctx context.Context, api *moov.APIClient) (time.Time, error) {
	downloads, resp, err := api.OFACApi.GetLatestDownloads(ctx, &moov.GetLatestDownloadsOpts{
		Limit: optional.NewInt32(1),
	})
	if err != nil {
		return time.Time{}, err
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return time.Time{}, fmt.Errorf("download error (stats code: %d): %v", resp.StatusCode, err)
	}
	if len(downloads) == 0 {
		return time.Time{}, errors.New("empty downloads response")
	}
	return downloads[0].Timestamp, nil
}
