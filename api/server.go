package api

import (
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	MasterQuery *db.MsQ
	Router      *gin.Engine
}

func NewServer(masterQ *db.MsQ) *Server {
	server := &Server{
		MasterQuery: masterQ,
	}

	router := gin.Default()

	router.POST("/accounts", server.createAccount)

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
