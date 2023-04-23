package api

import (
	"os"
	"testing"
	"time"

	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, msQ db.MsQ) *Server {
	config := utils.Config{
		SymmetricKey:  utils.RdmString(32),
		TokenDuration: time.Minute,
	}
	server, err := NewServer(config, msQ)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
