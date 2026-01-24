package main

import (
	"net/http"
	"sort"

	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")
	var (
		dbChirps []database.Chirp
		err      error
	)
	if authorIdString != "" {
		authorID, err := uuid.Parse(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Author ID", err)
			return
		}

		dbChirps, err = cfg.queries.GetChirpByAuthor(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.queries.GetChirps(r.Context())
	}

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
	sortMethod := r.URL.Query().Get("sort")
	if (sortMethod != "") && (sortMethod == "desc") {
		sort.Slice(
			resChirps,
			func(i, j int) bool { return resChirps[i].CreatedAt.After(resChirps[j].CreatedAt) },
		)
	}
	respondWithJson(w, http.StatusOK, resChirps)
}

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
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
	resChirp := resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJson(w, http.StatusOK, resChirp)
}
