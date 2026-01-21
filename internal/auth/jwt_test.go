package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwtTokenValidationValid(t *testing.T) {
	userId, _ := uuid.NewRandom()
	tokenSecret := "thehorseisonthestable"
	expiresIn, _ := time.ParseDuration("72h")
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Token generation failed: %v", err)
	}
	tokenUserId, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}
	if tokenUserId != userId {
		t.Fatal("UserID do not match")
	}
}

func TestJwtTokenValidationExpired(t *testing.T) {
	userId, _ := uuid.NewRandom()
	tokenSecret := "thehorseisonthestable"
	expiresIn, _ := time.ParseDuration("1ns")
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Token generation failed: %v", err)
	}
	time.Sleep(time.Millisecond)
	tokenUserId, err := ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil (userID=%v)", tokenUserId)
	}
}

func TestJwtTokenValidationWrongSecret(t *testing.T) {
	userId, _ := uuid.NewRandom()
	tokenSecret := "thehorseisonthestable"
	tokenSecretWrong := "!thehorseisonthestable"
	expiresIn, _ := time.ParseDuration("72h")
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Token generation failed: %v", err)
	}
	time.Sleep(time.Millisecond)
	tokenUserId, err := ValidateJWT(tokenString, tokenSecretWrong)
	if err == nil {
		t.Fatalf("expected error for wrong secret, got nil (userID=%v)", tokenUserId)
	}
}

func TestGetBearerToken(t *testing.T) {
	request := http.Request{}
	request.Header = map[string][]string{
		"Authorization": {"Bearer thehorseisonthestable"},
	}
	token, err := GetBearerToken(request.Header)
	if err != nil {
		t.Fatalf("GetBearerToken should not fail with error: %v", err)
	}
	if token != "thehorseisonthestable" {
		t.Fatalf("Token should be thehorseisonthestable but is: %s", token)
	}
}

func TestMakeRefreshToken(t *testing.T) {
	_, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("Refresh Token should be generated without errors")

	}
}









