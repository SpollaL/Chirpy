package main

import (
	"encoding/json"
	"net/http"

	"github.com/SpollaL/Chirpy/internal/auth"
	"github.com/SpollaL/Chirpy/internal/database"
)

type reqUserUpdate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error extracting token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret_key)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	userUpdate := reqUserUpdate{}
	err = decoder.Decode(&userUpdate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse request body", err)
	}
	hashedPassword, err := auth.HashPassword(userUpdate.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}
	dbUserUpdated, err := cfg.queries.UpdateUser(
		r.Context(),
		database.UpdateUserParams{
			Email:          userUpdate.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could not update users email and password",
			err,
		)
	}
	jsonUser := resUser{
		ID:        dbUserUpdated.ID,
		CreatedAt: dbUserUpdated.CreatedAt,
		UpdatedAt: dbUserUpdated.UpdatedAt,
		Email:     dbUserUpdated.Email,
	}
	respondWithJson(w, http.StatusOK, jsonUser)
}
