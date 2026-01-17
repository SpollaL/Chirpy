package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
	}
	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		log.Fatalf("Could not delete all users: %v", err)
	}
	cfg.fileserverHits.Store(0)
}

