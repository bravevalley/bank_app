package token

import (
	"testing"
	"time"

	"github.com/dassyareg/bank_app/utils"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJwToken_Happy(t *testing.T) {
	
	Jwtoken, err := NewJwToken(utils.RandomEmail(32))
	require.NoError(t, err)
	require.NotEmpty(t, Jwtoken)

	username := utils.RandomName()
	duration := time.Second * 3

	createdAt := time.Now()
	expiredAt := time.Now().Add(duration)

	Token, err := Jwtoken.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, Token)

	payload, err := Jwtoken.VerifyToken(Token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	
	require.NotEmpty(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}


func TestJwToken_ExpiredToken(t *testing.T) {
	
	Jwtoken, err := NewJwToken(utils.RandomEmail(32))
	require.NoError(t, err)
	require.NotEmpty(t, Jwtoken)

	username := utils.RandomName()

	Token, err := Jwtoken.GenerateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, Token)

	payload, err := Jwtoken.VerifyToken(Token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())

	require.Nil(t, payload)
	
}

func TestJwToken_ShortSecretKey(t *testing.T) {
	
	_, err := NewJwToken(utils.RdmString(31))
	require.Error(t, err)
	require.EqualError(t, err, "Invalid Secretkey: too short")
	
}

func TestJwToken_InvalidToken(t *testing.T) {
	payload, err := NewPayload(utils.RandomName(), time.Minute * 1)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	
	
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	Token, err := jwtoken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, Token)


	tokenizer, err := NewJwToken(utils.RdmString(32))
	require.NoError(t, err)

	NewPayload, err := tokenizer.VerifyToken(Token)
	require.Error(t, err)
	require.Empty(t, NewPayload)
	require.EqualError(t, err, ErrInvalidToken.Error())

	
	
}

