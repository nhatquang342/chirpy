package auth_test

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/nhatquang342/chirpy/internal/auth"
)

func TestJWTCreateAndValidate(t *testing.T) {
	secret := "supersecretkey"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotUID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	if gotUID != userID {
		t.Fatalf("expected %v, got %v", userID, gotUID)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "supersecretkey"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, -time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestJWTRongSecret(t *testing.T) {
	secret := "supersecretkey"
	wrong := "wrongkey"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = auth.ValidateJWT(token, wrong)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}
