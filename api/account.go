package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountParams struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount is used to create a new bank account
func (server *Server) createAccount(gc *gin.Context) {
	var createdAcc createAccountParams

	if err := gc.ShouldBindJSON(&createdAcc); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	clientData := gc.MustGet(AuthPayload).(*token.Payload)

	account, err := server.MasterQuery.CreateAccount(gc, db.CreateAccountParams{
		Name:     clientData.Username,
		Balance:  0,
		Currency: createdAcc.Currency,
	})

	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				gc.JSON(http.StatusForbidden, errorRes(err))
			}
			return
		}
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	gc.IndentedJSON(http.StatusCreated, account)
}

type getAccountByIDParams struct {
	ID int `uri:"id" binding:"required,min=0"`
}

func (server *Server) getAccountByID(gc *gin.Context) {
	var accountID getAccountByIDParams

	if err := gc.ShouldBindUri(&accountID); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	account, err := server.MasterQuery.GetAccount(gc, int64(accountID.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			gc.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	clientData := gc.MustGet(AuthPayload).(*token.Payload)
	if account.Name != clientData.Username {
		err := errors.New("User unauthorized to view data")
		gc.JSON(http.StatusUnauthorized, errorRes(err))
		return
	}


	gc.IndentedJSON(http.StatusOK, account)

}

type listAccountsParams struct {
	PageNumber int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(gc *gin.Context) {
	var listAcc listAccountsParams

	if err := gc.ShouldBindQuery(&listAcc); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	clientData := gc.MustGet(AuthPayload).(*token.Payload)
	xaccounts, err := server.MasterQuery.ListAccount(gc, db.ListAccountParams{
		Name: clientData.Username,
		Limit:  listAcc.PageSize,
		Offset: (listAcc.PageNumber - 1) * listAcc.PageSize,
	})

	if err != nil {
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	gc.IndentedJSON(http.StatusOK, xaccounts)
}
