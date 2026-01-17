package main

import (
	"encoding/json"
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

func (cfg *apiConfig) HandleChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := reqChirp{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Unable to decode json body"}`))
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Chirp is too long. Should be less than 140 characters"}`))
		return
	}

	chirp.Body = ReplaceProfane(chirp.Body)
	dbChirp, err := cfg.queries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{Body: chirp.Body, UserID: chirp.UserId},
	)
	resChirp := resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	res, err := json.Marshal(resChirp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Unable to marshal response body"}`))
		return
	}

	w.WriteHeader(201)
	w.Write(res)
}
