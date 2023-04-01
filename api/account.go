package api

import (
	"net/http"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountParams struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD NGN"`
}

func (server *Server) createAccount(gc *gin.Context) {
	var createdAcc createAccountParams

	if err := gc.ShouldBindJSON(&createdAcc); err != nil {
		gc.JSON(http.StatusBadRequest, errorRes(err))
		return
	}

	account, err := server.MasterQuery.CreateAccount(gc, db.CreateAccountParams{
		Name:     createdAcc.Name,
		Balance:  0,
		Currency: createdAcc.Currency,
	})

	if err != nil {
		gc.JSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	gc.IndentedJSON(http.StatusCreated, account)

}
