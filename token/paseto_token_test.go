package token

import (
	"testing"
	"time"

	"github.com/dassyareg/bank_app/utils"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPaseto_Happy(t *testing.T) {

	PaseToken, err := NewPaseToken(utils.RdmString(32))
	require.NoError(t, err)
	require.NotEmpty(t, PaseToken)

	username := utils.RandomName()
	duration := time.Second * 3

	createdAt := time.Now()
	expiredAt := time.Now().Add(duration)

	Token, err := PaseToken.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, Token)

	payload, err := PaseToken.VerifyToken(Token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPaseto_ExpiredToken(t *testing.T) {

	PaseToken, err := NewPaseToken(utils.RdmString(32))
	require.NoError(t, err)
	require.NotEmpty(t, PaseToken)

	username := utils.RandomName()

	Token, err := PaseToken.GenerateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, Token)

	payload, err := PaseToken.VerifyToken(Token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())

	require.Nil(t, payload)

}

func TestPaseto_ShortSecretKey(t *testing.T) {

	_, err := NewJwToken(utils.RdmString(31))
	require.Error(t, err)
	require.EqualError(t, err, "Invalid Secretkey: too short")

}

func TestPaseto_InvalidToken(t *testing.T) {
	payload, err := NewPayload(utils.RandomName(), time.Minute*1)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	PaseToken, err := paseto.NewV2().Encrypt([]byte("YELLOW SUBMARINE, BLACK WIZARDRY"), payload, nil)
	require.NoError(t, err)

	tokenizer, err := NewPaseToken(utils.RdmString(32))
	require.NoError(t, err)

	NewPayload, err := tokenizer.VerifyToken(PaseToken)
	require.Error(t, err)
	require.Empty(t, NewPayload)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
