package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockdb "github.com/dassyareg/bank_app/db/mocks"
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/token"
	"github.com/dassyareg/bank_app/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferTranx(t *testing.T) {
	// Create a randomUser
	userA := randomAccount(randomUser().Username) 
	userB := randomAccount(randomUser().Username)

	
	currency := utils.USD
	amount := int(utils.RandomAmount())
	
	userA.Currency = currency
	userB.Currency = currency

	testCases := []struct {
		Name	string
		Arg TransferTranxParams
		Stubs func(m *mockdb.MockMsQ)
		Expect func(t *testing.T, res *httptest.ResponseRecorder)
		Token func (t *testing.T, req *http.Request, tokenMaker token.TokenMaker)
		
	}{
		{
			Name: "Happy Case",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(2).Return(userA, nil)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Return(transferReturn(userA, userB, amount), nil)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Sender unauthorized",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(userB, nil)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0).Return(transferReturn(userA, userB, amount), nil)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Sender does not exist",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(db.Account{}, sql.ErrNoRows)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0).Return(transferReturn(userA, userB, amount), nil)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Receiver does not exist",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(userA, nil)

				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(db.Account{}, sql.ErrNoRows)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0).Return(transferReturn(userA, userB, amount), nil)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "SQL server downtime",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(db.Account{}, sql.ErrConnDone)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0).Return(transferReturn(userA, userB, amount), nil)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "SQL server downtime transfer",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(userA, nil)

				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(1).Return(userB, nil)

				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(1).Return(db.SuccessfulTransferResult{}, sql.ErrConnDone)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "unsupported currency",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: "ONE",
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userA.AccNumber) ).Times(0).Return(userA, nil)

				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userB.AccNumber)).Times(0).Return(userB, nil)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Invalid currency",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userA.AccNumber) ).Times(1).Return(db.Account{
					AccNumber: userA.AccNumber,
					Name: userA.Name,
					Balance: utils.RandomAmount(),
					Currency: "ONE",
				}, nil)

				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userB.AccNumber)).Times(1).Return(userB, nil)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Invalid currency",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userA.AccNumber) ).Times(1).Return(userA, nil)

				m.EXPECT().GetAccount(gomock.Any(), gomock.Eq(userB.AccNumber)).Times(1).Return(db.Account{
					AccNumber: userB.AccNumber,
					Name: userB.Name,
					Balance: utils.RandomAmount(),
					Currency: "ONE",
				}, nil)
	
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
		{
			Name: "Invalid input",
			Arg: TransferTranxParams{
				Sender: userA.AccNumber,
				Receiver: userB.AccNumber,
				Amount: -amount,
				Currency: currency,
			},
			Stubs: func(m *mockdb.MockMsQ) {
				m.EXPECT().GetAccount(gomock.Any(), mockdb.AccMatcher(userB.AccNumber, userA.AccNumber)).Times(0)
				m.EXPECT().ExecTransferTx(gomock.Any(), gomock.Eq(db.TransferProcessParams{
					Debit: userA.AccNumber,
					Credit: userB.AccNumber,
					Amount: int64(amount),
				})).Times(0)

			},
			Expect: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, res.Code)
			},
			Token: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, userA.Name, time.Minute, tokenMaker, req, AuthType)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			// Create the mock controller and interface
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mockdb.NewMockMsQ(ctrl)
			tC.Stubs(m)

			server := NewTestServer(t, m)

			url := fmt.Sprint("/transfers")

			res := httptest.NewRecorder()

			requestPayload, err := json.Marshal(&tC.Arg)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(requestPayload)))
			require.NoError(t, err)

			tC.Token(t, req, server.TokenMaker)

			server.Router.ServeHTTP(res, req)

			tC.Expect(t, res)
			
			
		})
	}
}

func transferReturn(sender db.Account, receiver db.Account, amt int ) db.SuccessfulTransferResult {
	return db.SuccessfulTransferResult{
		Transfer: db.Transfer{
			ID: utils.RandomAmount(),
			Amount: int64(amt),
			Debit: sender.AccNumber,
			Credit: receiver.AccNumber,
		},
		SenderAcc: sender,
		ReceiverAcc: receiver,
		SenderTransaction: db.Transaction{
			ID: utils.RandomAmount(),
			AccNumber: sender.AccNumber,
			Amount: -int64(amt),
		},
		ReceiverTransaction: db.Transaction{
			ID: utils.RandomAmount(),
			AccNumber: receiver.AccNumber,
			Amount: int64(amt),
		},
	}
}