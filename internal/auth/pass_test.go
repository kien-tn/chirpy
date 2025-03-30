package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword(%q) returned error: %v", password, err)
	}
	if hash == "" {
		t.Fatalf("HashPassword(%q) returned empty hash", password)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword(%q) returned error: %v", password, err)
	}
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Fatalf("CheckPasswordHash(%q, %q) returned error: %v", hash, password, err)
	}
	err = CheckPasswordHash(hash, "wrongpassword")
	if err == nil {
		t.Fatalf("CheckPasswordHash(%q, %q) did not return error", hash, "wrongpassword")
	}
}
