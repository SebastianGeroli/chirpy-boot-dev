package auth

import (
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
