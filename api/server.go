package api

import (
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	MasterQuery db.MsQ
	Router      *gin.Engine
}

func NewServer(masterQ db.MsQ) *Server {
	server := &Server{
		MasterQuery: masterQ,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountByID)
	router.GET("/accounts", server.listAccounts)
	router.POST("/transfers", server.TransferTranx)

	server.Router = router
	return server
}

// Start start the http server and listens for request on the address provided
func (server *Server) Start(addr string) error {
	return server.Router.Run(addr)
}

// Error response JSONfy
func errorRes(err error) gin.H {
	return gin.H{"error": err.Error()}
}
