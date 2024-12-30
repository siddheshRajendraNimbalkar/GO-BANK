package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(key string) (*PasetoMaker, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 bytes")
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}, nil
}

// CreateToken creates and signs a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, Payload, error) {
	// Create the payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", *payload, fmt.Errorf("failed to create payload: %w", err)
	}

	// Encrypt the payload to create the token
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", *payload, fmt.Errorf("failed to encrypt token: %w", err)
	}

	return token, *payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
