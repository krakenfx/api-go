package helper

import (
	"github.com/google/uuid"
)

// UUID creates a new random UUID and returns it as a string or panics
func UUID() string {
	return uuid.NewString()
}
