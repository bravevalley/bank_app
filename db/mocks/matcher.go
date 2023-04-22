package mockdb

import (
	"fmt"
	"reflect"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
)

type UserTestMatcher struct {
	Arg      db.CreateUserParams
	Password string
}

func (u UserTestMatcher) Matches(x interface{}) bool {
	// Assert the input implement the type
	UserInput, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// Verify the inputed password with the converted
	err := utils.VerifyPassword(u.Password, UserInput.HashedPassword)
	if err != nil {
		return false
	}

	// Set the value of the inputed password to the expected password
	u.Arg.HashedPassword = UserInput.HashedPassword

	// Check if the values input is the one expected
	return reflect.DeepEqual(u.Arg, UserInput)
}

func (u UserTestMatcher) String() string {
	return fmt.Sprintf("Matches %v with %v", u.Arg.HashedPassword, u.Password)
}

func MatchUserInput(argument db.CreateUserParams, password string) gomock.Matcher {
	return UserTestMatcher{
		Arg:      argument,
		Password: password,
	}
}
