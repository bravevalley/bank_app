package api

import (
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/token"
	"github.com/dassyareg/bank_app/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	TokenMaker  token.TokenMaker
	Config      utils.Config
	MasterQuery db.MsQ
	Router      *gin.Engine
}

func NewServer(config utils.Config, masterQ db.MsQ) (*Server, error) {
	Token, err := token.NewPaseToken(config.SymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		TokenMaker:  Token,
		MasterQuery: masterQ,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/users", server.addUser)
	router.POST("/users/login", server.loginUser)

	authRouteGroup := router.Group("/").Use(AuthMiddleWare(server.TokenMaker))

	authRouteGroup.POST("/accounts", server.createAccount)
	authRouteGroup.GET("/accounts/:id", server.getAccountByID)
	authRouteGroup.GET("/accounts", server.listAccounts)
	authRouteGroup.POST("/transfers", server.TransferTranx)

	server.Router = router
	return server, nil
}

// Start start the http server and listens for request on the address provided
func (server *Server) Start(addr string) error {
	return server.Router.Run(addr)
}

// Error response JSONfy
func errorRes(err error) gin.H {
	return gin.H{"error": err.Error()}
}
