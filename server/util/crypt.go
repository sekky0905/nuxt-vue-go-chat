package util

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns hashed password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), errors.Wrap(err, "failed to generate from password")
}

// CheckHashOfPassword checks whether given hashed is the value of hashed password or not.
func CheckHashOfPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
