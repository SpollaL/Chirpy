package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.queries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	resChirps := []resChirp{}
	for _, dbChirp := range dbChirps {
		resChirp := resChirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		resChirps = append(resChirps, resChirp)
	}
	respondWithJson(w, http.StatusOK, resChirps)
}

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirp", err)
	}
	dbChirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not retrieve chirp", err)
	}
	resChirp := resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJson(w, http.StatusOK, resChirp)
}
