package auth

import (
	"testing"
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
