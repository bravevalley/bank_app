package token

import "time"

// TokenMaker implements the behaviour for JOSE token package
type TokenMaker interface {

	// GenerateToken creates a token from the username and duration
	GenerateToken(username string, duration time.Duration) (string, error)

	// VerifyToken Verifies and authenticate a token
	VerifyToken(token *string) (*Payload, error)
}