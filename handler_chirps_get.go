package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	resChirps := []resChirp{}
	dbChirps, err := cfg.queries.GetChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not get chirps from database: %v", err)
		return
	}
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
	res, err := json.Marshal(resChirps)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Unable to marshal response body"}`))
		return
	}
	w.WriteHeader(200)
	w.Write(res)
}

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not parse ID %v: %v", chirpID, err)
		return
	}
	dbChirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	resChirp := resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Could not get chirp with ID %v: %v", chirpID, err)
		return
	}
	res, err := json.Marshal(resChirp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could marshal chirp with ID %v: %v", chirpID, err)
		return
	}
	w.WriteHeader(200)
	w.Write(res)
}
