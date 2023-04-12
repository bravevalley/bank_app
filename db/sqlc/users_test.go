package db

import (
	"context"
	"testing"
	"time"

	"github.com/dassyareg/bank_app/utils"
	"github.com/stretchr/testify/require"
)

func randomUser() User {
	return User{
		Username: utils.RandomName(),
		HashedPassword: utils.RandomName(),
		FullName: utils.RandomName(),
		Email: utils.RandomName(),
	}
}



// TestCreateUser tests the CreateUser method
func CreateAUser(t *testing.T) User {
	newUser := randomUser()

	createdUser, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username: newUser.Username,
		HashedPassword: newUser.HashedPassword,
		FullName: newUser.FullName,
		Email: newUser.Email,
	})

	require.NoError(t, err)
	require.NotEmpty(t, createdUser)

	require.Equal(t, newUser.Username, createdUser.Username)
	require.Equal(t, newUser.HashedPassword, createdUser.HashedPassword)
	require.Equal(t, newUser.FullName, createdUser.FullName)
	require.Equal(t, newUser.Email, createdUser.Email)

	require.NotEmpty(t, createdUser.CreatedAt)
	require.NotEmpty(t, createdUser.PasswordLastChanged)

	return createdUser
}

func TestCreateUser(t *testing.T) {
	_ = CreateAUser(t)
}

func TestGetUser(t *testing.T) {
	wantUser := CreateAUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), wantUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	require.Equal(t, wantUser.Username, gotUser.Username)
	require.Equal(t, wantUser.FullName, gotUser.FullName)
	require.Equal(t, wantUser.Email, gotUser.Email)
	require.Equal(t, wantUser.HashedPassword, gotUser.HashedPassword)

	require.NotEmpty(t, gotUser.CreatedAt)
	require.NotEmpty(t, gotUser.PasswordLastChanged)

	require.WithinDuration(t, wantUser.CreatedAt, gotUser.CreatedAt, time.Second)
	require.WithinDuration(t, wantUser.PasswordLastChanged, gotUser.PasswordLastChanged, time.Second)

}