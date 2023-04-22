package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrExpiredToken is returned when a token has expired
	ErrExpiredToken = errors.New("Token Exired")

	// ErrInvalidToken is returned when a token fails authentication
	ErrInvalidToken = errors.New("Invalid Token Detected")
)

// Payload is the struct containing the payload for the authentication
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdat"`
	ExpiredAt time.Time `json:"expiredat"`
}

// NewPayload create a new Payload from a username and duration, returns a
// pointer to the payload
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id := uuid.Must(uuid.NewRandom())

	return &Payload{
		ID:        id,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// Valid is used to validate the Token by using the Payload
func (pl *Payload) Valid() error {
	if time.Now().After(pl.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
