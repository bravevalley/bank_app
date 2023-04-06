package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/gin-gonic/gin"
)

// TransferTranxParams is the model for the possible api request
type TransferTranxParams struct {
	Sender   int64  `json:"sender" binding:"required"`
	Receiver int64  `json:"receiver" binding:"required"`
	Amount   int    `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,oneof=USD NGN EUR"`
}

// TransferTranx is the handler for the transfer api endpoint
func (server *Server) TransferTranx(gc *gin.Context) {
	var transferTransaction TransferTranxParams

	if err := gc.BindJSON(&transferTransaction); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}
	
	if ok := server.testCurrency(gc, transferTransaction.Sender, transferTransaction.Currency); !ok {
		gc.JSON(http.StatusBadRequest, errorRes(fmt.Errorf("Currency of Account %d doesn't match the submitted currency [%v]", transferTransaction.Sender, transferTransaction.Currency)))
		return
	}
	if ok := server.testCurrency(gc, transferTransaction.Receiver, transferTransaction.Currency); !ok {
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
func (server *Server) testCurrency(gc *gin.Context, accNumber int64, currency string) bool {
	account, err := server.MasterQuery.GetAccount(gc, accNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			gc.JSON(http.StatusNotFound, errorRes(err))
			return false
		}
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return false
	}

	if account.Currency != currency {
		return false
	}
	return true
}
