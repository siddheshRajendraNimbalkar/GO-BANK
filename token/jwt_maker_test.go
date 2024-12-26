package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJwt(t *testing.T) {
	// Define the symmetric key for the token
	key := "12345678901234567890123456789012"

	// Create a new PASETO token maker
	maker, err := NewPasetoMaker(key)
	require.NoError(t, err)

	// Create a token
	username := "example_user"
	duration := time.Hour * 24

	// Generate a token using the maker
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	fmt.Println("Generated Token:", token)

	// Verify the token
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	// Validate the payload
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt, time.Second)

	fmt.Printf("Token verified! Payload: %+v\n", payload)
}
