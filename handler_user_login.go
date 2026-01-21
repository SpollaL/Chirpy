package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SpollaL/Chirpy/internal/auth"
	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/google/uuid"
)

type reqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type resLogin struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqStruct := &reqLogin{}
	err := decoder.Decode(reqStruct)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	dbUser, err := cfg.queries.GetUser(r.Context(), reqStruct.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user", err)
		return
	}
	match, err := auth.CheckPasswordHash(reqStruct.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password", err)
		return
	}
	const defaultExpiresInSeconds = 60
	token, err := auth.MakeJWT(
		dbUser.ID,
		cfg.secret_key,
		time.Duration(defaultExpiresInSeconds)*time.Second,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate JWT token", err)
		return
	}
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate Refresh token", err)
		return
	}
	const defaultExpiresInDays = 60
	_, err = cfg.queries.CreateToken(
		r.Context(),
		database.CreateTokenParams{
			Token:     refresh_token,
			UserID:    dbUser.ID,
			ExpiresAt: time.Now().AddDate(0, 0, defaultExpiresInDays),
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate Refresh token", err)
		return
	}
	jsonUser := resLogin{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refresh_token,
	}
	respondWithJson(w, http.StatusOK, jsonUser)
}
