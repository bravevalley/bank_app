package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// MinSecretKeyLen is the minimum length of characters that be used for
// Secret key
const MinSecretKeyLen = 32

// JwToken is the struct containing the secret key
type JwToken struct {
	Secretkey string 
}

// NewJwToken returns a new JWToken struct after testing its length against
// the minimum length
func NewJwToken(secretKey string) (TokenMaker, error) {
	if len(secretKey) < MinSecretKeyLen {
		return &JwToken{}, fmt.Errorf("Invalid Secretkey: too short")
	}

	return &JwToken{
		Secretkey: secretKey,
	}, nil
}

// GenerateToken creates a JWT token from the payload and signs it with the 
// JwToken struct secret key
func (jwtoken *JwToken) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("Cant create token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(jwtoken.Secretkey))
}

// VerifyToken verifies the authenticity of the token and parse the
// payload into the Payload struct
func (jwttoken *JwToken) VerifyToken(token string) (*Payload, error) {
	
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwttoken.Secretkey), nil
	}

	Token, err := jwt.ParseWithClaims(token , &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := Token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
