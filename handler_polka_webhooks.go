package main

import (
	"encoding/json"
	"net/http"

	"github.com/SpollaL/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type reqPolkaData struct {
	UserId string `json:"user_id"`
}

type reqPolka struct {
	Event string       `json:"event"`
	Data  reqPolkaData `json:"data"`
}

func (cfg *apiConfig) HandlePolkaWebHooks(w http.ResponseWriter, r *http.Request) {
  apiKey, err := auth.GetAPIKey(r.Header)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Couldn't retrieve api token", err)
    return
  }
  if apiKey != cfg.polka_key {
    respondWithError(w, http.StatusUnauthorized, "Couldn't verify api token", err)
    return
  }
	decoder := json.NewDecoder(r.Body)
	req := reqPolka{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse request body", err)
		return
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userId, err := uuid.Parse(req.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user id", err)
		return
	}
	_, err = cfg.queries.UpgradeUser(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
