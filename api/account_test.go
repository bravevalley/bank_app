package api

import (
	"fmt"
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

	msQ.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.AccNumber)).Times(1).Return(account, nil)

	url := fmt.Sprintf("/accounts/%d", account.AccNumber)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	res := httptest.NewRecorder()

	server := NewServer(msQ)
	server.Router.ServeHTTP(res, req)

	require.Equal(t, res.Code, http.StatusOK)
	
}


func randomAccount() db.Account {
	return db.Account{
		AccNumber: utils.RdmNumbBtwRange(2, 200),
		Name: utils.RandomName(),
		Balance: utils.RandomAmount(),
		Currency: utils.RdnCurr(),
	}
}