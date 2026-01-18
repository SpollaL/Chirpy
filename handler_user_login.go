package main

import (
	"encoding/json"
	"net/http"

	"github.com/SpollaL/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqStruct := &reqUser{}
	err := decoder.Decode(reqStruct)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
	}

	dbUser, err := cfg.queries.GetUser(r.Context(), reqStruct.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user", err)
	}
	match, err := auth.CheckPasswordHash(reqStruct.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password", err)
	}
	jsonUser := resUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJson(w, http.StatusOK, jsonUser)
}
