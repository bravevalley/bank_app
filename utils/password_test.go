package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
		testCases := []struct {
		Name	string
		Password string
		WantErr   error
		TestFunc func(t *testing.T, pwd string, err error)
		
	}{
		{
			Name: "Hash Password",
			Password: rdmString(6),
			WantErr: fmt.Errorf(""),
			TestFunc: func(t *testing.T, pwd string, wantErr error) {
				hashedPwd, err := HashPassword(pwd)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPwd)
			},
		},
		{
			Name: "Hash Uniqueness",
			Password: rdmString(6),
			WantErr: fmt.Errorf(""),
			TestFunc: func(t *testing.T, pwd string, wantErr error) {
				hashedPwd, err := HashPassword(pwd)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPwd)

				hashedPwd2, _ := HashPassword(pwd)
				require.NotEqual(t, hashedPwd, hashedPwd2)
				
			},
		},
		{
			Name: "Cant Hash",
			Password: rdmString(10000),
			WantErr: fmt.Errorf("Can't Hash Password"),
			TestFunc: func(t *testing.T, pwd string, wantErr error) {
				hashedPwd, err := HashPassword(pwd)
				require.Error(t, err)
				require.Empty(t, hashedPwd)
				require.EqualError(t, err,wantErr.Error())		
			},
		},
		{
			Name: "Verify Password: Right Input",
			Password: rdmString(6),
			WantErr: fmt.Errorf(""),
			TestFunc: func(t *testing.T, pwd string, wantErr error) {
				hashedPwd, err := HashPassword(pwd)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPwd)

				err = VerifyPassword(pwd, hashedPwd)
				require.NoError(t, err)	
			},
		},
		{
			Name: "Verify Password: Wrong Input",
			Password: rdmString(6),
			WantErr: fmt.Errorf(""),
			TestFunc: func(t *testing.T, pwd string, wantErr error) {
				hashedPwd, err := HashPassword(pwd)
				require.NoError(t, err)
				require.NotEmpty(t, hashedPwd)

				err = VerifyPassword("test", hashedPwd)
				require.Error(t, err)
				require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.TestFunc(t, tC.Password, tC.WantErr)
		})
	}


}

