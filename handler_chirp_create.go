package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type reqChirp struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type resChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func ReplaceProfane(s string) string {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	for word := range strings.SplitSeq(s, " ") {
		for _, pword := range profane_words {
			if strings.ToLower(word) == pword {
				splits := strings.Split(s, word)
				s = strings.Join(splits, "****")
			}
		}
	}
	return s
}

func validateChirp(body string) (string, error) {
	const maxChirpLenght = 140
	if len(body) > maxChirpLenght {
		return "", errors.New("Chirp is too long")
	}
	cleaned_body := ReplaceProfane(body)
	return cleaned_body, nil
}

func (cfg *apiConfig) HandleChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := reqChirp{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

  cleaned, err := validateChirp(chirp.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbChirp, err := cfg.queries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{Body: cleaned, UserID: chirp.UserId},
	)
	resChirp := resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJson(w, http.StatusCreated, resChirp)
}
