package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.Nil, errors.New("empty UUID string")
	}

	uuidParsed, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	return uuidParsed, nil
}
