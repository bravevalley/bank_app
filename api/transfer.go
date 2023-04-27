package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/token"
	"github.com/gin-gonic/gin"
)

// TransferTranxParams is the model for the possible api request
type TransferTranxParams struct {
	Sender   int64  `json:"sender" binding:"required"`
	Receiver int64  `json:"receiver" binding:"required"`
	Amount   int    `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

// TransferTranx is the handler for the transfer api endpoint
func (server *Server) TransferTranx(gc *gin.Context) {
	var transferTransaction TransferTranxParams

	if err := gc.BindJSON(&transferTransaction); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	clientData := gc.MustGet(AuthPayload).(*token.Payload)

	senderAccount, ok := server.checkAcc(gc, transferTransaction.Sender)
	if !ok {
		return 
	}

	if senderAccount.Name != clientData.Username {
		err := errors.New("User unauthorized to make transfers")
		gc.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}

	receiverAccount, ok := server.checkAcc(gc, transferTransaction.Receiver)
	if !ok {
		return 
	}

	if ok := server.testCurrency(gc, senderAccount, transferTransaction.Currency); !ok {
		gc.JSON(http.StatusBadRequest, errorRes(fmt.Errorf("Currency of Account %d doesn't match the submitted currency [%v]", transferTransaction.Sender, transferTransaction.Currency)))
		return
	}
	if ok := server.testCurrency(gc, receiverAccount, transferTransaction.Currency); !ok {
		gc.JSON(http.StatusBadRequest, errorRes(fmt.Errorf("Currency of Account %d doesn't match the submitted currency [%v]", transferTransaction.Receiver, transferTransaction.Currency)))
		return
	}

	result, err := server.MasterQuery.ExecTransferTx(gc, db.TransferProcessParams{
		Debit:  transferTransaction.Sender,
		Credit: transferTransaction.Receiver,
		Amount: int64(transferTransaction.Amount),
	})

	if err != nil {
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	gc.IndentedJSON(http.StatusCreated, result)
}

// testCurrenty tests if client has an account with the currency submitted. Returns false when there is any error or if the currency do not match
func (server *Server) testCurrency(gc *gin.Context, acc db.Account, currency string) bool {
	
	if acc.Currency != currency {
		return false
	}
	return true
}

func (server *Server) checkAcc(gc *gin.Context, accNumber int64) (db.Account, bool) {
	account, err := server.MasterQuery.GetAccount(gc, accNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			gc.JSON(http.StatusNotFound, errorRes(err))
			return db.Account{}, false
		}
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return db.Account{}, false
	}

	return account, true
}
