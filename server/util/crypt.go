package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns hashed password.
func HashPassword(password string) (string, error) {
	cost := 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckHashOfPassword checks whether given hashed is the value of hashed password or not.
func CheckHashOfPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
