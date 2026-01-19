package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	method := jwt.SigningMethodHS256
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(method, claims)
	stringifiedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return stringifiedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, err
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer_token := headers.Get("Authorization")
	if bearer_token == "" {
		return "", errors.New("No Authorization header found in request")
	}
	token, match := strings.CutPrefix(bearer_token, "Bearer ")
	if !match {
		return "", errors.New("Couldn't find Bearer section in Authorization header")
	}
	return token, nil
}
