package main

import (
	"net/http"
)

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithJson(w, http.StatusForbidden, "Reset is only allowed in development")
	}
	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to restard the database", err)
		return
	}
	cfg.fileserverHits.Store(0)
	respondWithJson(w, http.StatusOK, "Hits reset to 0 and dataset to initial state")
}

