package api

import (
	"database/sql"
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
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func randomUser() db.User {
	return db.User{
		Username:       utils.RandomName(),
		HashedPassword: utils.RdmString(7),
		FullName:       utils.RandomName(),
		Email:          utils.RandomEmail(5),
	}
}

func TestAddUser(t *testing.T) {
	var errCode *pq.Error

	errCode = &pq.Error{
		Code: "23505",
	}

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
		{
			Name: "Wrong data input",
			Arg: CreateUserArgs{
				Username: User.Username,
				Password: Password,
				FullName: User.FullName,
				Email:    "testemail_com",
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().
					CreateUser(gomock.Any(), mockdb.MatchUserInput(db.CreateUserParams{
						Username:       User.Username,
						HashedPassword: Password,
						FullName:       User.FullName,
						Email:          User.Email,
					}, Password)).Times(0)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			Name: "Unhashable password",
			Arg: CreateUserArgs{
				Username: User.Username,
				Password: utils.RdmString(1000),
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
					}, Password)).Times(0)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
		},
		{
			Name: "SQL Server Error",
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
					Return(db.User{}, sql.ErrConnDone)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
		},
		{
			Name: "Constraint Error",
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
					Return(db.User{}, errCode)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, res.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.Stub(mockDB)

			server := NewTestServer(t, mockDB)
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

func TestUserLogin(t *testing.T) {

	

	// Create the mock controller and interface
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockdb.NewMockMsQ(ctrl)

	// Create a randomUser
	user := randomUser()

	// Hash the password for verfification
	hashed, err := utils.HashPassword(user.HashedPassword)
	require.NoError(t, err)

	testCases := []struct {
		Name     string
		Arg      UserLoginArgs
		Stub     func(m *mockdb.MockMsQ)
		Expected func(t *testing.T, res *httptest.ResponseRecorder)
	}{
		{
			Name: "Happy Case",
			Arg: UserLoginArgs{
				Username: user.Username,
				Password: user.HashedPassword,
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{
					Username: user.Username,
					HashedPassword: hashed,
					FullName: user.Username,
					Email: user.Email,
				}, nil)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, res.Code)
			},
		},
		{
			Name: "Wrong Password",
			Arg: UserLoginArgs{
				Username: user.Username,
				Password: user.HashedPassword,
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{
					Username: user.Username,
					HashedPassword: user.HashedPassword,
					FullName: user.Username,
					Email: user.Email,
				}, nil)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, res.Code)
			},
		},
		{
			Name: "User does not exist",
			Arg: UserLoginArgs{
				Username: utils.RdmString(3),
				Password: user.HashedPassword,
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
		{
			Name: "SQL server down",
			Arg: UserLoginArgs{
				Username: user.Username,
				Password: user.HashedPassword,
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
		},
		{
			Name: "Invalid JSON Input",
			Arg: UserLoginArgs{
				Username: user.Username,
				Password: utils.RdmString(3),
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetUser(gomock.Any(), user.Username).Times(0)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {

			tC.Stub(m)

			server := NewTestServer(t, m)

			url := fmt.Sprint("/users/login")

			res := httptest.NewRecorder()

			requestPayload, err := json.Marshal(&tC.Arg)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(requestPayload)))
			require.NoError(t, err)

			server.Router.ServeHTTP(res, req)

			tC.Expected(t, res)
		})
	}
}
