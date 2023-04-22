package api

import (
	"net/http"
	"time"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserArgs struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type ResponseUser struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) addUser(gc *gin.Context) {
	var createdUser CreateUserArgs

	if err := gc.ShouldBindJSON(&createdUser); err != nil {
		gc.IndentedJSON(http.StatusBadRequest, errorRes(err))
		return
	}

	hashedPwd, err := utils.HashPassword(createdUser.Password)
	if err != nil {
		gc.IndentedJSON(http.StatusInternalServerError, errorRes(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       createdUser.Username,
		HashedPassword: hashedPwd,
		FullName:       createdUser.FullName,
		Email:          createdUser.Email,
	}

	user, err := server.MasterQuery.CreateUser(gc, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation":
				gc.IndentedJSON(http.StatusForbidden, errorRes(err))
				return
			}
			gc.IndentedJSON(http.StatusInternalServerError, errorRes(err))
			return
		}
	}

	AddedUser := ResponseUser{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	gc.IndentedJSON(http.StatusCreated, AddedUser)

}
