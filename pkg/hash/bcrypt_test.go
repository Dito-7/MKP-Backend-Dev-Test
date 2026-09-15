package hash

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Errorf("expected password to match hash")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Errorf("expected wrong password to fail")
	}
	t.Logf("Hash for password123: %s", hash)
}
