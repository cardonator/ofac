// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestDownload__manualRefreshPath(t *testing.T) {
	if testing.Short() {
		return
	}

	searcher := &searcher{}
	repo := createTestDownloadRepository(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", manualRefreshPath, nil)
	logger := log.NewNopLogger()
	manualRefreshHandler(logger, searcher, repo)(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}
	var stats downloadStats
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Error(err)
	}
	if stats.SDNs == 0 {
		t.Errorf("stats.SDNs=%d but expected non-zero", stats.SDNs)
	}
	if stats.Alts == 0 {
		t.Errorf("stats.Alts=%d but expected non-zero", stats.Alts)
	}
	if stats.Addresses == 0 {
		t.Errorf("stats.Addresses=%d but expected non-zero", stats.Addresses)
	}
	if stats.DeniedPersons == 0 {
		t.Errorf("stats.DeniedPersons=%d but expected non-zero", stats.DeniedPersons)
	}
	if stats.SectoralSanctions == 0 {
		t.Errorf("stats.SectoralSanctions=%d but expected non-zero", stats.SectoralSanctions)
	}
}
