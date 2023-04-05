package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/dassyareg/bank_app/db/mocks"
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestCreateAccount(t *testing.T) {
	account := randomAccount()

	mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()

	msQ := mock.NewMockMsQ(mockCtrl)

	
	testCases := []struct {
		Name string
		AccNumber int64
		stub func (m *mock.MockMsQ)
		expected func (t *testing.T, response *httptest.ResponseRecorder) 
	
	}{
		{
			Name: "Ok",
			AccNumber: account.AccNumber,
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.AccNumber)).Times(1).Return(account, nil)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusOK)
				testResponseBody(t, res.Body, account)
			},
			
		},
		{
			Name: "Bad Request",
			AccNumber: -1,
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.AccNumber)).Times(0)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
			},
			
		},
		{
			Name: "No rows",
			AccNumber: account.AccNumber,
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.AccNumber)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			expected: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusNotFound)
			},
			
		},
		{
			Name: "Internal Server Error",
			AccNumber: account.AccNumber,
			stub: func(m *mock.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.AccNumber)).Times(1).Return(db.Account{}, sql.ErrConnDone)
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
		
			server := NewServer(msQ)
			server.Router.ServeHTTP(res, req)

			tC.expected(t, res)
		})
	}

}


func testResponseBody(t *testing.T, body *bytes.Buffer, acc db.Account)  {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAcc db.Account

	err = json.Unmarshal(data, &gotAcc)
	require.NoError(t, err)
	require.Equal(t, acc, gotAcc)
}

func randomAccount() db.Account {
	return db.Account{
		AccNumber: utils.RdmNumbBtwRange(2, 200),
		Name: utils.RandomName(),
		Balance: utils.RandomAmount(),
		Currency: utils.RdnCurr(),
	}
}