package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type reqUser struct {
	Email string `json:"email"`
}

type resUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) HandleUserCreation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqStruct := &reqUser{}
	err := decoder.Decode(reqStruct)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
	}
	dbUser, err := cfg.queries.CreateUser(r.Context(), reqStruct.Email)
	jsonUser := resUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJson(w, http.StatusCreated, jsonUser)
}
