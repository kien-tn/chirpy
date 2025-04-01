package auth

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "01234567890123456789012345678901"
	expiresIn := time.Minute * 5

	// Test MakeJWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	log.Printf("token: %v", token)
	if err != nil {
		t.Fatalf("Error creating JWT: %v", err)
	}

	// Test ValidateJWT
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if parsedUserID != userID {
		t.Fatalf("Expected user ID %s, got %s", userID, parsedUserID)
	}
}
