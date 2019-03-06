package util

import "github.com/google/uuid"

// UUID generates UUID.
func UUID() string {
	return uuid.New().String()
}
