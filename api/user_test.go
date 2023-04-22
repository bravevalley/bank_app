package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockdb "github.com/dassyareg/bank_app/db/mocks"
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomUser() db.User {
	return db.User{
		Username:       utils.RandomName(),
		HashedPassword: utils.RandomName(),
		FullName:       utils.RandomName(),
		Email:          utils.RandomEmail(5),
	}
}

func TestAddUser(t *testing.T) {
	Password := fmt.Sprint("123456789")
	User := randomUser()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mockdb.NewMockMsQ(mockCtrl)

	testCases := []struct {
		Name     string
		Arg      CreateUserArgs
		Stub     func(m *mockdb.MockMsQ)
		Expected func(t *testing.T, res *httptest.ResponseRecorder)
	}{
		{
			Name: "Status Created",
			Arg: CreateUserArgs{
				Username: User.Username,
				Password: Password,
				FullName: User.FullName,
				Email:    User.Email,
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().
					CreateUser(gomock.Any(), mockdb.MatchUserInput(db.CreateUserParams{
						Username:       User.Username,
						HashedPassword: Password,
						FullName:       User.FullName,
						Email:          User.Email,
					}, Password)).Times(1).
					Return(User, nil)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, res.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.Stub(mockDB)

			server := NewServer(mockDB)
			res := httptest.NewRecorder()

			body, err := json.Marshal(&tC.Arg)
			require.NoError(t, err)

			url := fmt.Sprint("/users")
			req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))

			server.Router.ServeHTTP(res, req)

			tC.Expected(t, res)

		})
	}
}
