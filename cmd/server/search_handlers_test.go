// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cardonator/ofac"

	"github.com/gorilla/mux"
)

func TestSearch__Address(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.89`) {
		t.Errorf("%#v", v)
	}

	var wrapper struct {
		Addresses []*ofac.Address `json:"addresses"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.Addresses[0].EntityID != "173" {
		t.Errorf("%#v", wrapper.Addresses[0])
	}

	// send an empty body and get an error
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/search?limit=1", nil)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("bogus status code: %d", w.Code)
	}
}

func TestSearch__AddressCountry(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?country=united+kingdom&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressMulti(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.945`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressProvidence(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&providence=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.96333`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressCity(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&city=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.96333`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressState(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&state=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.96333`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__NameAndAltName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?limit=1&q=Air+I", nil)

	s := &searcher{
		Alts:      altSearcher.Alts,
		SDNs:      sdnSearcher.SDNs,
		Addresses: addressSearcher.Addresses,
		DPs:       dplSearcher.DPs,
		SSIs:      ssiSearcher.SSIs,
		ELs:       elSearcher.ELs,
	}

	router := mux.NewRouter()
	addSearchRoutes(nil, router, s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	// read response body
	var wrapper struct {
		SDNs              []*ofac.SDN               `json:"SDNs"`
		AltNames          []*ofac.AlternateIdentity `json:"altNames"`
		Addresses         []*ofac.Address           `json:"addresses"`
		DeniedPersons     []*ofac.DPL               `json:"deniedPersons"`
		SectoralSanctions []*ofac.SSI               `json:"sectoralSanctions"`
		BISEntities       []*ofac.EL                `json:"bisEntities"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.SDNs[0].EntityID != "2676" {
		t.Errorf("%#v", wrapper.SDNs[0])
	}
	if wrapper.AltNames[0].EntityID != "4691" {
		t.Errorf("%#v", wrapper.AltNames[0].EntityID)
	}
	if wrapper.Addresses[0].EntityID != "735" {
		t.Errorf("%#v", wrapper.Addresses[0].EntityID)
	}
	if wrapper.DeniedPersons[0].StreetAddress != "P.O. BOX 28360" {
		t.Errorf("%#v", wrapper.DeniedPersons[0].StreetAddress)
	}
	if wrapper.SectoralSanctions[0].EntityID != "18736" {
		t.Errorf("%#v", wrapper.SectoralSanctions[0].EntityID)
	}
	if wrapper.BISEntities[0].Name != "Mohammad Jan Khan Mangal" {
		t.Errorf("%#v", wrapper.BISEntities[0].Name)
	}
}

func TestSearch__Name(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?name=AL+ZAWAHIRI&limit=1", nil)

	router := mux.NewRouter()
	combinedSearcher := &searcher{
		SDNs: sdnSearcher.SDNs,
		DPs:  dplSearcher.DPs,
		SSIs: ssiSearcher.SSIs,
		ELs:  elSearcher.ELs,
	}
	addSearchRoutes(nil, router, combinedSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.91`) {
		t.Error(v)
	}

	var wrapper struct {
		SDNs []*ofac.SDN `json:"SDNs"`
		DPs  []*ofac.DPL `json:"deniedPersons"`
		SSIs []*ofac.SSI `json:"sectoralSanctions"`
		ELs  []*ofac.EL  `json:"bisEntities"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if len(wrapper.SDNs) != 1 || len(wrapper.SSIs) != 1 || len(wrapper.DPs) != 1 {
		t.Fatalf("SDNs=%d SSIs=%d DPs=%d", len(wrapper.SDNs), len(wrapper.SSIs), len(wrapper.DPs))
	}
	if wrapper.SDNs[0].EntityID != "2676" {
		t.Errorf("%#v", wrapper.SDNs[0])
	}
	if wrapper.SSIs[0].EntityID != "18736" {
		t.Errorf("%#v", wrapper.SSIs[0])
	}
	if wrapper.DPs[0].Name != "AL NASER WINGS AIRLINES" {
		t.Errorf("%#v", wrapper.DPs[0])
	}
	if wrapper.ELs[0].Name != "Luqman Yasin Yunus Shgragi" {
		t.Errorf("%#v", wrapper.ELs[0])
	}
}

func TestSearch__AltName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?altName=sogo+KENKYUSHO&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(nil, router, &searcher{
		Alts: altSearcher.Alts,
	})
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.783`) {
		t.Error(v)
	}

	var wrapper struct {
		Alts []*ofac.AlternateIdentity `json:"altNames"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if len(wrapper.Alts) != 1 {
		t.Fatalf("Alts=%d", len(wrapper.Alts))
	}
	if wrapper.Alts[0].EntityID != "4691" {
		t.Errorf("%#v", wrapper.Alts[0])
	}
}
