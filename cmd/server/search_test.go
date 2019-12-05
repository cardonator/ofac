// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cardonator/ofac"

	"github.com/go-kit/kit/log"
)

var (
	addressSearcher = &searcher{
		Addresses: precomputeAddresses([]*ofac.Address{
			{
				EntityID:                    "173",
				AddressID:                   "129",
				Address:                     "Ibex House, The Minories",
				CityStateProvincePostalCode: "London EC3N 1DY",
				Country:                     "United Kingdom",
			},
			{
				EntityID:                    "735",
				AddressID:                   "447",
				Address:                     "Piarco Airport",
				CityStateProvincePostalCode: "Port au Prince",
				Country:                     "Haiti",
			},
		}),
	}
	altSearcher = &searcher{
		Alts: precomputeAlts([]*ofac.AlternateIdentity{
			{ // Real OFAC entry
				EntityID:      "559",
				AlternateID:   "481",
				AlternateType: "aka",
				AlternateName: "CIMEX",
			},
			{
				EntityID:      "4691",
				AlternateID:   "3887",
				AlternateType: "aka",
				AlternateName: "A.I.C. SOGO KENKYUSHO",
			},
		}),
	}
	sdnSearcher = &searcher{
		SDNs: precomputeSDNs([]*ofac.SDN{
			{
				EntityID: "2676",
				SDNName:  "AL ZAWAHIRI, Dr. Ayman",
				SDNType:  "individual",
				Program:  "SDGT] [SDT",
				Title:    "Operational and Military Leader of JIHAD GROUP",
				Remarks:  "DOB 19 Jun 1951; POB Giza, Egypt; Passport 1084010 (Egypt); alt. Passport 19820215; Operational and Military Leader of JIHAD GROUP.",
			},
			{
				EntityID: "2681",
				SDNName:  "HAWATMA, Nayif",
				SDNType:  "individual",
				Program:  "SDT",
				Title:    "Secretary General of DEMOCRATIC FRONT FOR THE LIBERATION OF PALESTINE - HAWATMEH FACTION",
				Remarks:  "DOB 1933; Secretary General of DEMOCRATIC FRONT FOR THE LIBERATION OF PALESTINE - HAWATMEH FACTION.",
			},
		}),
	}
	dplSearcher = &searcher{
		DPs: precomputeDPs([]*ofac.DPL{
			{
				Name:           "AL NASER WINGS AIRLINES",
				StreetAddress:  "P.O. BOX 28360",
				City:           "DUBAI",
				State:          "",
				Country:        "AE",
				PostalCode:     "",
				EffectiveDate:  "06/05/2019",
				ExpirationDate: "12/03/2019",
				StandardOrder:  "Y",
				LastUpdate:     "2019-06-12",
				Action:         "FR NOTICE ADDED, TDO RENEWAL, F.R. NOTICE ADDED, TDO RENEWAL ADDED, TDO RENEWAL ADDED, F.R. NOTICE ADDED",
				FRCitation:     "82 F.R. 61745 12/29/2017,  83F.R. 28801 6/21/2018, 84 F.R. 27233 6/12/2019",
			},
			{
				Name:           "PRESTON JOHN ENGEBRETSON",
				StreetAddress:  "12725 ROYAL DRIVE",
				City:           "STAFFORD",
				State:          "TX",
				Country:        "US",
				PostalCode:     "77477",
				EffectiveDate:  "01/24/2002",
				ExpirationDate: "01/24/2027",
				StandardOrder:  "Y",
				LastUpdate:     "2002-01-28",
				Action:         "STANDARD ORDER",
				FRCitation:     "67 F.R. 7354 2/19/02 66 F.R. 48998 9/25/01 62 F.R. 26471 5/14/97 62 F.R. 34688 6/27/97 62 F.R. 60063 11/6/97 63 F.R. 25817 5/11/98 63 F.R. 58707 11/2/98 64 F.R. 23049 4/29/99",
			},
		}),
	}
	ssiSearcher = &searcher{
		SSIs: precomputeSSIs([]*ofac.SSI{
			{
				EntityID:       "18782",
				Type:           "Entity",
				Programs:       []string{"SYRIA", "UKRAINE-EO13662"},
				Name:           "ROSOBORONEKSPORT OAO",
				Addresses:      []string{"27 Stromynka ul., Moscow, 107076, RU"},
				Remarks:        []string{"For more information on directives, please visit the following link: http://www.treasury.gov/resource-center/sanctions/Programs/Pages/ukraine.aspx#directives", "(Linked To: ROSTEC)"},
				AlternateNames: []string{"RUSSIAN DEFENSE EXPORT ROSOBORONEXPORT", "ROSOBORONEXPORT JSC", "ROSOBORONEKSPORT OJSC", "OJSC ROSOBORONEXPORT", "ROSOBORONEXPORT"},
				IDsOnRecord:    []string{"1117746521452, Registration ID", "56467052, Government Gazette Number", "7718852163, Tax ID No.", "Subject to Directive 3, Executive Order 13662 Directive Determination -", "www.roe.ru, Website"},
				SourceListURL:  "http://bit.ly/1QWTIfE",
				SourceInfoURL:  "http://bit.ly/1MLgou0",
			},
			{
				EntityID:       "18736",
				Type:           "Entity",
				Programs:       []string{"UKRAINE-EO13662"},
				Name:           "VTB SPECIALIZED DEPOSITORY, CJSC",
				Addresses:      []string{"35 Myasnitskaya Street, Moscow, 101000, RU"},
				Remarks:        []string{"For more information on directives, please visit the following link: http://www.treasury.gov/resource-center/sanctions/Programs/Pages/ukraine.aspx#directives", "(Linked To: ROSTEC)"},
				AlternateNames: []string{"CJS VTB SPECIALIZED DEPOSITORY"},
				IDsOnRecord:    []string{"1117746521452, Registration ID", "56467052, Government Gazette Number", "7718852163, Tax ID No.", "Subject to Directive 3, Executive Order 13662 Directive Determination -", "www.roe.ru, Website"},
				SourceListURL:  "http://bit.ly/1QWTIfE",
				SourceInfoURL:  "http://bit.ly/1MLgou0",
			},
		}),
	}
	elSearcher = &searcher{
		ELs: precomputeELs([]*ofac.EL{
			{
				Name:               "Mohammad Jan Khan Mangal",
				AlternateNames:     []string{"Air I"},
				Addresses:          []string{"Kolola Pushta, Charahi Gul-e-Surkh, Kabul, AF", "Maidan Sahr, Hetefaq Market, Paktiya, AF"},
				StartDate:          "11/13/19",
				LicenceRequirement: "For all items subject to the EAR (See ¬ß744.11 of the EAR). ",
				LicensePolicy:      "Presumption of denial.",
				FRNotice:           "81 FR 57451",
				SourceListURL:      "http://bit.ly/1L47xrV",
				SourceInfoURL:      "http://bit.ly/1L47xrV",
			},
			{
				Name:               "Luqman Yasin Yunus Shgragi",
				AlternateNames:     []string{"Lkemanasel Yosef", "Luqman Sehreci."},
				Addresses:          []string{"Savcili Mahalesi Turkmenler Caddesi No:2, Sahinbey, Gaziantep, TR", "Sanayi Mahalesi 60214 Nolu Caddesi No 11, SehitKamil, Gaziantep, TR"},
				StartDate:          "8/23/16",
				LicenceRequirement: "For all items subject to the EAR.  (See ¬ß744.11 of the EAR)",
				LicensePolicy:      "Presumption of denial.",
				FRNotice:           "81 FR 57451",
				SourceListURL:      "http://bit.ly/1L47xrV",
				SourceInfoURL:      "http://bit.ly/1L47xrV",
			},
		}),
	}
)

func TestJaroWrinkler(t *testing.T) {
	cases := []struct {
		s1, s2 string
		match  float64
	}{
		{"wei, zhao", "wei, Zhao", 0.950},
		{"WEI, Zhao", "WEI, Zhao", 1.0},
		// make sure jaroWrinkler is communative
		{"jane doe", "jan lahore", 0.69},
		{"jan lahore", "jane doe", 0.69},
		// example cases
		{"maduro moros, nicolas", "maduro moros, nicolas", 1.0},
		{"maduro moros, nicolas", "nicolas maduro", 0.512},
		{"nicolas maduro moros", "nicolás maduro", 0.855},
		{"nicolas, maduro moros", "nicolas maduro", 0.891},
		{"nicolas, maduro moros", "nicolás maduro", 0.881},
	}

	for _, v := range cases {
		// Only need to call chomp on s1, see jaroWrinkler doc
		eql(t, fmt.Sprintf("%s vs %s", v.s1, v.s2), jaroWrinkler(chomp(v.s1), v.s2), v.match)
	}
}

func eql(t *testing.T, desc string, x, y float64) {
	t.Helper()
	if math.Abs(x-y) > 0.01 {
		t.Errorf("%s: %.3f != %.3f", desc, x, y)
	}
}

func TestEql(t *testing.T) {
	eql(t, "", 0.1, 0.1)
	eql(t, "", 0.0001, 0.00002)
}

// TestSearch_precompute ensures we are trimming and UTF-8 normalizing strings
// as expected. This is needed since our datafiles are normalized for us.
func TestSearch_precompute(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"nicolás maduro", "nicolasmaduro"},
		{"Delcy Rodríguez", "delcyrodriguez"},
		{"Raúl Castro", "raulcastro"},
	}
	for i := range cases {
		guess := precompute(cases[i].input)
		if guess != cases[i].expected {
			t.Errorf("precompute(%q)=%q expected %q", cases[i].input, guess, cases[i].expected)
		}
	}
}

func TestSearch_reorderSDNName(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"Jane Doe", "Jane Doe"},                         // control
		{"Jane, Doe Other", "Jane, Doe Other"},           // made up name to make sure we don't clobber ,'s in the middle of a name
		{"FELIX B. MADURO S.A.", "FELIX B. MADURO S.A."}, // keep .'s in a name
		{"MADURO MOROS, Nicolas", "Nicolas MADURO MOROS"},
		{"IBRAHIM, Sadr", "Sadr IBRAHIM"},
	}
	for i := range cases {
		guess := reorderSDNName(cases[i].input, "individual")
		if guess != cases[i].expected {
			t.Errorf("reorderSDNName(%q)=%q expected %q", cases[i].input, guess, cases[i].expected)
		}
	}
}

// TestSearch_liveData will download the real OFAC data and run searches against the corpus.
// This test is designed to tweak match percents and results.
func TestSearch_liveData(t *testing.T) {
	if testing.Short() {
		return
	}
	searcher := &searcher{
		logger: log.NewNopLogger(),
	}
	if stats, err := searcher.refreshData(); err != nil {
		t.Fatal(err)
	} else {
		searcher.logger.Log("liveData", fmt.Sprintf("stats: %#v", stats))
	}

	cases := []struct {
		name  string
		match float64 // top match %
	}{
		{"Nicolas MADURO", 0.944},
	}
	for i := range cases {
		sdns := searcher.TopSDNs(1, cases[i].name)
		if len(sdns) == 0 {
			t.Errorf("name=%q got no results", cases[i].name)
		}
		eql(t, fmt.Sprintf("%q (SDN=%s) matches %q ", cases[i].name, sdns[0].EntityID, sdns[0].name), sdns[0].match, cases[i].match)
	}
}

func TestSearch__topAddressesAddress(t *testing.T) {
	it := topAddressesAddress("needle")(&Address{address: "needleee"})

	eql(t, "topAddressesAddress", it.weight, 0.95)
	if add, ok := it.value.(*Address); !ok || add.address != "needleee" {
		t.Errorf("got %#v", add)
	}
}

func TestSearch__topAddressesCountry(t *testing.T) {
	it := topAddressesAddress("needle")(&Address{address: "needleee"})

	eql(t, "topAddressesCountry", it.weight, 0.95)
	if add, ok := it.value.(*Address); !ok || add.address != "needleee" {
		t.Errorf("got %#v", add)
	}
}

func TestSearch__multiAddressCompare(t *testing.T) {
	it := multiAddressCompare(
		topAddressesAddress("needle"),
		topAddressesCountry("other"),
	)(&Address{address: "needlee", country: "other"})

	eql(t, "multiAddressCompare", it.weight, 0.9857)
	if add, ok := it.value.(*Address); !ok || add.address != "needlee" || add.country != "other" {
		t.Errorf("got %#v", add)
	}
}

func TestSearch__extractSearchLimit(t *testing.T) {
	// Too high, fallback to hard max
	req := httptest.NewRequest("GET", "/?limit=1000", nil)
	if limit := extractSearchLimit(req); limit != hardResultsLimit {
		t.Errorf("got limit of %d", limit)
	}

	// No limit, use default
	req = httptest.NewRequest("GET", "/", nil)
	if limit := extractSearchLimit(req); limit != softResultsLimit {
		t.Errorf("got limit of %d", limit)
	}

	// Between soft and hard max
	req = httptest.NewRequest("GET", "/?limit=25", nil)
	if limit := extractSearchLimit(req); limit != 25 {
		t.Errorf("got limit of %d", limit)
	}

	// Lower than soft max
	req = httptest.NewRequest("GET", "/?limit=1", nil)
	if limit := extractSearchLimit(req); limit != 1 {
		t.Errorf("got limit of %d", limit)
	}
}

func TestSearch__addressSearchRequest(t *testing.T) {
	u, _ := url.Parse("https://moov.io/search?address=add&city=new+york&state=ny&providence=prov&zip=44433&country=usa")
	req := readAddressSearchRequest(u)
	if req.Address != "add" {
		t.Errorf("req.Address=%s", req.Address)
	}
	if req.City != "new york" {
		t.Errorf("req.City=%s", req.City)
	}
	if req.State != "ny" {
		t.Errorf("req.State=%s", req.State)
	}
	if req.Providence != "prov" {
		t.Errorf("req.Providence=%s", req.Providence)
	}
	if req.Zip != "44433" {
		t.Errorf("req.Zip=%s", req.Zip)
	}
	if req.Country != "usa" {
		t.Errorf("req.Country=%s", req.Country)
	}
	if req.empty() {
		t.Error("req is not empty")
	}

	req = addressSearchRequest{}
	if !req.empty() {
		t.Error("req is empty now")
	}
	req.Address = "1600 1st St"
	if req.empty() {
		t.Error("req is not empty now")
	}
}

func TestSearch__FindAddresses(t *testing.T) {
	addresses := addressSearcher.FindAddresses(1, "173")
	if v := len(addresses); v != 1 {
		t.Fatalf("len(addresses)=%d", v)
	}
	if addresses[0].EntityID != "173" {
		t.Errorf("got %#v", addresses[0])
	}
}

func TestSearch__TopAddresses(t *testing.T) {
	addresses := addressSearcher.TopAddresses(1, "Piarco Air")
	if len(addresses) == 0 {
		t.Fatal("empty Addresses")
	}
	if addresses[0].Address.EntityID != "735" {
		t.Errorf("%#v", addresses[0].Address)
	}
}

func TestSearch__TopAddressFn(t *testing.T) {
	addresses := addressSearcher.TopAddressesFn(1, topAddressesCountry("United Kingdom"))
	if len(addresses) == 0 {
		t.Fatal("empty Addresses")
	}
	if addresses[0].Address.EntityID != "173" {
		t.Errorf("%#v", addresses[0].Address)
	}
}

func TestSearch__FindAlts(t *testing.T) {
	alts := altSearcher.FindAlts(1, "559")
	if v := len(alts); v != 1 {
		t.Fatalf("len(alts)=%d", v)
	}
	if alts[0].EntityID != "559" {
		t.Errorf("got %#v", alts[0])
	}
}

func TestSearch__TopSdnAlts(t *testing.T) {
	alts := altSearcher.TopAltNames(1, "SOGO KENKYUSHO")
	if len(alts) == 0 {
		t.Fatal("empty AltNames")
	}
	if alts[0].AlternateIdentity.EntityID != "4691" {
		t.Errorf("%#v", alts[0].AlternateIdentity)
	}
}

func TestSearch__FindSDN(t *testing.T) {
	sdn := sdnSearcher.FindSDN("2676")
	if sdn == nil {
		t.Fatal("nil SDN")
	}
	if sdn.EntityID != "2676" {
		t.Errorf("got %#v", sdn)
	}
}

func TestSearch__TopSDNs(t *testing.T) {
	sdns := sdnSearcher.TopSDNs(1, "AL ZAWAHIRI")
	if len(sdns) == 0 {
		t.Fatal("empty SDNs")
	}
	if sdns[0].EntityID != "2676" {
		t.Errorf("%#v", sdns[0].SDN)
	}
}

func TestSearch__TopDPs(t *testing.T) {
	dps := dplSearcher.TopDPs(1, "NASER AIRLINES")
	if len(dps) == 0 {
		t.Fatal("empty DPs")
	}
	// DPL doesn't have any entity IDs. Comparing expected address components instead
	if dps[0].DeniedPerson.StreetAddress != "P.O. BOX 28360" || dps[0].DeniedPerson.City != "DUBAI" {
		t.Errorf("%#v", dps[0].DeniedPerson)
	}
}

func TestSearcher_TopSSIs(t *testing.T) {
	ssis := ssiSearcher.TopSSIs(1, "ROSOBORONEKSPORT")
	if len(ssis) == 0 {
		t.Fatal("empty SSIs")
	}
	if ssis[0].SectoralSanction.EntityID != "18782" {
		t.Errorf("%#v", ssis[0].SectoralSanction)
	}
}

func TestSearcher_TopELs(t *testing.T) {
	els := elSearcher.TopELs(1, "Mohammad")
	if len(els) == 0 {
		t.Fatal("empty ELs")
	}
	if els[0].Entity.Name != "Mohammad Jan Khan Mangal" {
		t.Errorf("%#v", els[0].Entity)
	}
}
