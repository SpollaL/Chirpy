package main

import (
	"errors"
	"net/http"

	"github.com/SpollaL/Chirpy/internal/auth"
	"github.com/google/uuid"
)


func (cfg *apiConfig) HandleChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret_key)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirp", err)
		return
	}
	dbChirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not retrieve chirp", err)
		return
	}
	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "Unauthorized request to delete token", errors.New("Could not delete chirp"))
		return
	}
	err = cfg.queries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
