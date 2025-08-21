package user

import (
	"github.com/google/uuid"
)

func GenerateAuthToken() string {
	return uuid.New().String()
}
