package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(p), err
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
