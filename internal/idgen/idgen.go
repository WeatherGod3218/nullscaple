package idgen

import (
	"github.com/google/uuid"
)

func GenerateNewId() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	stringForm := uuid.String()
	return stringForm, nil
}
