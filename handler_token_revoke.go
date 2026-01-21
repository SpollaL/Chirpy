package main

import (
	"net/http"

	"github.com/SpollaL/Chirpy/internal/auth"
)


func (cfg *apiConfig) HandleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh token", err)
		return
	}
	err = cfg.queries.RevokeToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
