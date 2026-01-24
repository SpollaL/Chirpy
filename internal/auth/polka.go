package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	api_token := headers.Get("Authorization")
	if api_token == "" {
		return "", errors.New("No Authorization header found in request")
	}
	token, match := strings.CutPrefix(api_token, "ApiKey ")
	if !match {
		return "", errors.New("Couldn't find ApiKey section in Authorization header")
	}
	return token, nil
}
