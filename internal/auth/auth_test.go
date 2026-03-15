package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Errorf("expected no error, got %v:", err)
	}
	if hash == "" {
		t.Error("Expected hash to not be empty")
	}
	if hash == "password123" {
		t.Error("hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v:", err)
	}

	match, err := CheckPasswordHash("password123", hash)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !match {
		t.Error("expected password to match hash")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v:", err)
	}
	match, err :=  CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if match {
		t.Error("expected password to not match hash")
	}
}

func TestMakeJWTValid(t *testing.T) {
	secret := "testsecret"
	userID := uuid.New()
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Errorf("expecting no error, got %v", err)
	}
	if token == "" {
		t.Error("expecting token to not be empty")
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "testsecret"
	userID := uuid.New()
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("unable to create token: %v", err)
	}
	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("expecting no error, got %v", err)
	}
	if userID != gotID {
		t.Errorf("expecting %v, got %v", userID, gotID)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "correctsecret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Error("expecting error with wrong secret, got nil")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("unable to create token: %v", err)
	}
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("expecting err for expired token, got nil")
	}
}
