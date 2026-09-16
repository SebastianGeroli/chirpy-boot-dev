package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "some-secret"

	token, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if gotID != userID {
		t.Errorf("expected userID %v, got %v", userID, gotID)
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "some-secret"

	token, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "some-secret"
	wrongSecret := "wrong-secret"

	token, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("expected error for token signed with wrong secret, got nil")
	}
}

func TestValidateBearer(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer skdjaskldjlksad")
	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}
	if token != "skdjaskldjlksad" {
		t.Errorf("expected token: %v got: %v", "skdjaskldjlksad", token)
	}
}

func TestValidateEmptyBearerHeader(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Error("expected error for empty Authorization header, got nil")
	}
}

func TestValidateBearerMissingPrefix(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "skdjaskldjlksad")
	_, err := GetBearerToken(header)
	if err == nil {
		t.Error("expected error for Authorization header missing Bearer prefix, got nil")
	}
}
