package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRandomEmail(t *testing.T) {
	Email := RandomEmail(6)
	validate := validator.New()
	err := validate.Var(Email, "required,email")
	require.NotEmpty(t, Email)
	require.NoError(t, err)
}
