package main

import (
	"net/http"
	"time"

	"github.com/SpollaL/Chirpy/internal/auth"
)

type resToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh token", err)
		return
	}
	dbRefreshToken, err := cfg.queries.GetToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token", err)
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) || dbRefreshToken.RevokedAt.Valid {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Refresh Token has expired or has been revoked",
			err,
		)
		return
	}
	const defaultExpiresInSeconds = 60
	token, err := auth.MakeJWT(
		dbRefreshToken.UserID,
		cfg.secret_key,
		time.Duration(defaultExpiresInSeconds)*time.Second,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate JWT token", err)
		return
	}
	respondWithJson(w, http.StatusOK, resToken{Token: token})
}
