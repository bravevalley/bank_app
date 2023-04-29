package mockdb

import (
	"fmt"
	"reflect"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
)


// UserTestMatcher is implements the matcher Interface
// so it can be used to match the password input with the
// hashed password
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


type AccNumberMatcher struct {
	SenderAccNo int64
	ReceiverAccNo int64
	Input int64
}

func (a AccNumberMatcher) Matches(x interface{}) bool {
	inputAcc, ok := x.(int64)
	if !ok {
		return false
	}
	a.Input = inputAcc

	var tempVar int64
	
	switch inputAcc {
	case a.SenderAccNo:
		tempVar = a.SenderAccNo
	case a.ReceiverAccNo:
		tempVar = a.ReceiverAccNo

	}
	return reflect.DeepEqual(tempVar, inputAcc)
}


func (a AccNumberMatcher) String() string {
	return fmt.Sprintf("Matches %d with %d or %d", a.Input, a.SenderAccNo, a.ReceiverAccNo)
}

func AccMatcher(creditAcc, debitAcc int64) gomock.Matcher {
	return AccNumberMatcher{
		SenderAccNo: debitAcc,
		ReceiverAccNo: creditAcc,
	}
}