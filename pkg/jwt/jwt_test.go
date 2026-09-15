package jwt

import (
	"testing"
	"time"
)

func TestJWT_GenerateAndValidate(t *testing.T) {
	secret := "my-secret-key-12345"
	userID := "user-uuid-1"
	email := "user@mkp.com"
	role := "ADMIN"

	token, err := GenerateToken(userID, email, role, secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Errorf("expected error with wrong secret, got nil")
	}
}
