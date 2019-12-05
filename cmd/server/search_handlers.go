// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	moovhttp "github.com/moov-io/base/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/search").HandlerFunc(search(logger, searcher))
}

type addressSearchRequest struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Providence string `json:"providence"`
	Zip        string `json:"zip"`
	Country    string `json:"country"`
}

func (req addressSearchRequest) empty() bool {
	return req.Address == "" && req.City == "" && req.State == "" &&
		req.Providence == "" && req.Zip == "" && req.Country == ""
}

func readAddressSearchRequest(u *url.URL) addressSearchRequest {
	return addressSearchRequest{
		Address:    strings.ToLower(strings.TrimSpace(u.Query().Get("address"))),
		City:       strings.ToLower(strings.TrimSpace(u.Query().Get("city"))),
		State:      strings.ToLower(strings.TrimSpace(u.Query().Get("state"))),
		Providence: strings.ToLower(strings.TrimSpace(u.Query().Get("providence"))),
		Zip:        strings.ToLower(strings.TrimSpace(u.Query().Get("zip"))),
		Country:    strings.ToLower(strings.TrimSpace(u.Query().Get("country"))),
	}
}

func search(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)

		// Search over all fields
		if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
			if logger != nil {
				logger.Log("search", fmt.Sprintf("searching all names and address for %s", q))
			}
			searchViaQ(logger, searcher, q)(w, r)
			return
		}

		// Search by Name
		if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
			if logger != nil {
				logger.Log("search", fmt.Sprintf("searching SDN names for %s", name))
			}
			searchByName(logger, searcher, name)(w, r)
			return
		}

		// Search by Alt Name
		if alt := strings.TrimSpace(r.URL.Query().Get("altName")); alt != "" {
			if logger != nil {
				logger.Log("search", fmt.Sprintf("searching SDN alt names for %s", alt))
			}
			searchByAltName(logger, searcher, alt)(w, r)
			return
		}

		// Search Addresses
		if req := readAddressSearchRequest(r.URL); !req.empty() {
			if logger != nil {
				logger.Log("search", fmt.Sprintf("searching address for %#v", req))
			}
			searchByAddress(logger, searcher, req)(w, r)
			return
		}

		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

type searchResponse struct {
	SDNs              []SDN     `json:"SDNs"`
	AltNames          []Alt     `json:"altNames"`
	Addresses         []Address `json:"addresses"`
	DeniedPersons     []DP      `json:"deniedPersons"`
	SectoralSanctions []SSI     `json:"sectoralSanctions"`
}

func searchByAddress(logger log.Logger, searcher *searcher, req addressSearchRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if req.empty() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var resp searchResponse
		limit := extractSearchLimit(r)

		var compares []func(*Address) *item
		if req.Address != "" {
			compares = append(compares, topAddressesAddress(req.Address))
		}

		if req.City != "" {
			compares = append(compares, topAddressesCityState(req.City))
		}
		if req.State != "" {
			compares = append(compares, topAddressesCityState(req.State))
		}
		if req.Providence != "" {
			compares = append(compares, topAddressesCityState(req.Providence))
		}
		if req.Zip != "" {
			compares = append(compares, topAddressesCityState(req.Zip))
		}
		if req.Country != "" {
			compares = append(compares, topAddressesCountry(req.Country))
		}

		// Perform our ranking across all accumulated compare functions
		//
		// TODO(adam): Is there something in the (SDN?) files which signal to block an entire country? (i.e. Needing to block Iran all together)
		// https://www.treasury.gov/resource-center/sanctions/CivPen/Documents/20190327_decker_settlement.pdf
		resp.Addresses = searcher.TopAddressesFn(limit, multiAddressCompare(compares...))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func searchViaQ(logger log.Logger, searcher *searcher, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name = strings.TrimSpace(name)
		if name == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)

		response := &searchResponse{
			SDNs:              searcher.TopSDNs(limit, name),
			AltNames:          searcher.TopAltNames(limit, name),
			Addresses:         searcher.TopAddresses(limit, name),
			DeniedPersons:     searcher.TopDPs(limit, name),
			SectoralSanctions: searcher.TopSSIs(limit, name),
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func searchByName(logger log.Logger, searcher *searcher, nameSlug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameSlug = strings.TrimSpace(nameSlug)
		if nameSlug == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)
		sdns := searcher.TopSDNs(limit, nameSlug)
		dps := searcher.TopDPs(limit, nameSlug)
		ssis := searcher.TopSSIs(limit, nameSlug)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&searchResponse{
			SDNs:              sdns,
			DeniedPersons:     dps,
			SectoralSanctions: ssis,
		})
		if err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func searchByAltName(logger log.Logger, searcher *searcher, altSlug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		altSlug = strings.TrimSpace(altSlug)
		if altSlug == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)
		alts := searcher.TopAltNames(limit, altSlug)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&searchResponse{
			AltNames: alts,
		})
		if err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}
