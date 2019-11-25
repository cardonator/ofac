// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
)

const manualRefreshPath = "/ofac/refresh"

// manualRefreshHandler will register an endpoint on the admin server OFAC data refresh endpoint
func manualRefreshHandler(logger log.Logger, searcher *searcher, downloadRepo downloadRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Log("main", "admin: refreshing sanctions lists")
		}
		if stats, err := searcher.refreshData(); err != nil {
			if logger != nil {
				logger.Log("main", fmt.Sprintf("ERROR: admin: problem refreshing sanctions lists: %v", err))
			}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if logger != nil {
				logger.Log("main", fmt.Sprintf("admin: finished sacntions lists refresh - Addresses=%d AltNames=%d SDNs=%d DeniedPersons=%d SecotralSanctions=%d ELs=%d",
					stats.Addresses, stats.Alts, stats.SDNs, stats.DeniedPersons, stats.SectoralSanctions, stats.BISEntities))
			}
			downloadRepo.recordStats(stats)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stats)
		}
	}
}
