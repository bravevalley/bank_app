package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const MinSecretKeyLen = 32

type JwToken struct {
	Secretkey string 
}

func NewJwToken(secretKey string) (TokenMaker, error) {
	if len(secretKey) < MinSecretKeyLen {
		return &JwToken{}, fmt.Errorf("Invalid Secretkey: too short")
	}

	return &JwToken{
		Secretkey: secretKey,
	}, nil
}


func (jwtoken *JwToken) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("Cant create token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

		return token.SignedString(jwtoken.Secretkey)
	}

	
func (jwt *JwToken) VerifyToken(token *string) (*Payload, error)
