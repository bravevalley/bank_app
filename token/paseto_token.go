package token

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PaseToken struct {
	Version      paseto.V2
	SymmetricKey []byte
}

func NewPaseToken(SymmetricKey string) (TokenMaker, error) {
	if len(SymmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrInvalidToken
	}

	return &PaseToken{
		Version:      *paseto.NewV2(),
		SymmetricKey: []byte(SymmetricKey),
	}, nil
}

func (pst *PaseToken) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := pst.Version.Encrypt(pst.SymmetricKey, payload, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (pst *PaseToken) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := pst.Version.Decrypt(token, pst.SymmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
