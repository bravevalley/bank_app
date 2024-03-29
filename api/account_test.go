package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mock "github.com/dassyareg/bank_app/db/mocks"
	mockdb "github.com/dassyareg/bank_app/db/mocks"
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/token"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountByID(t *testing.T) {
	account := randomAccount(utils.RandomName())

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msQ := mock.NewMockMsQ(mockCtrl)

	testCases := []struct {
		Name      string
		AccNumber int64
		Setup func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker)
		stub      func(m *mock.MockMsQ)
		expected  func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			Name:      "Ok",
			AccNumber: account.AccNumber,
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, account.Name, time.Minute, tokenMaker, req, AuthType)
			},
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.AccNumber)).Times(1).Return(account, nil)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusOK)
				testResponseBody(t, res.Body, account)
			},
		},
		{
			Name:      "Bad Request",
			AccNumber: -1,
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, account.Name, time.Minute, tokenMaker, req, AuthType)
			},
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.AccNumber)).Times(0)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
			},
		},
		{
			Name:      "No rows",
			AccNumber: account.AccNumber,
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, account.Name, time.Minute, tokenMaker, req, AuthType)
			},
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.AccNumber)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusNotFound)
			},
		},
		{
			Name:      "Internal Server Error",
			AccNumber: account.AccNumber,
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, account.Name, time.Minute, tokenMaker, req, AuthType)
			},
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.AccNumber)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusInternalServerError)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.stub(msQ)

			url := fmt.Sprintf("/accounts/%d", tC.AccNumber)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()

			server := NewTestServer(t, msQ)


			tC.Setup(t, req,server.TokenMaker)
			server.Router.ServeHTTP(res, req)

			tC.expected(t, res)
		})
	}

}

func testResponseBody(t *testing.T, body *bytes.Buffer, acc db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAcc db.Account

	err = json.Unmarshal(data, &gotAcc)
	require.NoError(t, err)
	require.Equal(t, acc, gotAcc)
}

func randomAccount(owner string) db.Account {
	return db.Account{
		AccNumber: utils.RdmNumbBtwRange(2, 200),
		Name:      owner,
		Balance:   utils.RandomAmount(),
		Currency:  utils.RdnCurr(),
	}
}

func TestCreateAccount(t *testing.T) {
	userName := utils.RandomName()
	

	account := db.Account{
		Name:      userName,
		AccNumber: utils.RdmNumbBtwRange(1, 500),
		Balance:   0,
		Currency:  utils.RdnCurr(),
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msQ := mockdb.NewMockMsQ(mockCtrl)

	testCases := []struct {
		Name        string
		AccountInfo createAccountParams
		Setup func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker)
		Stub        func(m *mock.MockMsQ)
		Expected    func(t *testing.T, res *httptest.ResponseRecorder)
	}{
		{
			Name: "Ok",
			AccountInfo: createAccountParams{
				Currency: account.Currency,
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userName, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mock.MockMsQ) {
				m.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{	
					Name: userName,
					Balance:  account.Balance,
					Currency: account.Currency,
				})).Times(1).Return(account, nil)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusCreated)
				testResponseBody(t, res.Body, account)
			},
		},
		{
			Name: "Internal Server Error",
			AccountInfo: createAccountParams{
				Currency: account.Currency,
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userName, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mock.MockMsQ) {
				m.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Name: userName,
					Balance:  account.Balance,
					Currency: account.Currency,
				})).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusInternalServerError)
			},
		},
		{
			Name: "Bad Request",
			AccountInfo: createAccountParams{
				Currency: "POP",
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userName, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mock.MockMsQ) {
				m.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Name: userName,
					Balance:  account.Balance,
					Currency: account.Currency,
				})).Times(0)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.Stub(msQ)

			server := NewTestServer(t, msQ)
			url := fmt.Sprint("/accounts")

			reqBody, err := json.Marshal(&tC.AccountInfo)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(string(reqBody)))

			res := httptest.NewRecorder()

			tC.Setup(t, req, server.TokenMaker)

			server.Router.ServeHTTP(res, req)

			tC.Expected(t, res)
		})
	}
}

func TestListAccounts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	msQ := mock.NewMockMsQ(mockCtrl)

	var xacc []db.Account
	user := randomUser()

	fakes := 5

	for i := 0; i < fakes; i++ {
		xacc = append(xacc, randomAccount(user.Username))
	}

	testCases := []struct {
		Name     string
		Input    listAccountsParams
		Setup func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker)
		Stub     func(m *mock.MockMsQ)
		Expected func(t *testing.T, res *httptest.ResponseRecorder)
	}{
		{
			Name: "List account",
			Input: listAccountsParams{
				PageNumber: 1,
				PageSize:   5,
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, user.Username, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().ListAccount(gomock.Any(), gomock.Eq(db.ListAccountParams{
					Name: user.Username,
					Limit:  5,
					Offset: 0,
				})).Times(1).Return(xacc, nil)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusOK)

				var gotxAcc []db.Account

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				err = json.Unmarshal(data, &gotxAcc)
				require.NoError(t, err)

				require.Equal(t, xacc, gotxAcc)
			},
		},
		{
			Name: "Internal Server Error",
			Input: listAccountsParams{
				PageNumber: 1,
				PageSize:   5,
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, user.Username, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().ListAccount(gomock.Any(), gomock.Eq(db.ListAccountParams{
					Name: user.Username,
					Limit:  5,
					Offset: 0,
				})).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusInternalServerError)
			},
		},
		{
			Name: "Bad Request",
			Input: listAccountsParams{
				PageNumber: 1,
				PageSize:   200,
			},
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, user.Username, time.Minute, tokenMaker, req, AuthType)
			},
			Stub: func(m *mockdb.MockMsQ) {
				m.EXPECT().ListAccount(gomock.Any(), gomock.Eq(db.ListAccountParams{
					Name: user.Username,
					Limit:  5,
					Offset: 0,
				})).Times(0)
			},
			Expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			tC.Stub(msQ)

			server := NewTestServer(t, msQ)

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tC.Input.PageNumber, tC.Input.PageSize)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()

			tC.Setup(t, req, server.TokenMaker)
			server.Router.ServeHTTP(res, req)

			tC.Expected(t, res)

		})
	}
}
