package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SpollaL/Chirpy/internal/auth"
	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type reqUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	hashedPassword, err := auth.HashPassword(reqStruct.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
	}
	dbUser, err := cfg.queries.CreateUser(
		r.Context(),
		database.CreateUserParams{Email: reqStruct.Email, HashedPassword: hashedPassword},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
	}
	jsonUser := resUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJson(w, http.StatusCreated, jsonUser)
}
