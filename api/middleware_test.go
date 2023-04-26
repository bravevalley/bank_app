package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dassyareg/bank_app/token"
	"github.com/dassyareg/bank_app/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func createTokenReq(t *testing.T, username string, duration time.Duration, tokenMaker token.TokenMaker)  string {
	token, err :=  tokenMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	return token
}

func TestAuthMiddleWare(t *testing.T) {
	tests := []struct {
		Name  string
		Setup func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker)
		CheckResponse func(t *testing.T, res httptest.ResponseRecorder)
	}{
		{
			Name: "Ok",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				Token := createTokenReq(t, utils.RandomName(), time.Minute, tokenMaker)
	
				AuthValue := fmt.Sprintf("%s %s", AuthType, Token)
	
				req.Header.Set(AuthHeaderField, AuthValue)
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusOK)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			url := fmt.Sprint("/auth")
			
			server.Router.GET("/auth", AuthMiddleWare(server.TokenMaker), func(gc *gin.Context) {
				gc.JSON(http.StatusOK, gin.H{
					"hello": "world",
				})
			})

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// TODO: Modify Request with Token and AUtorization Header
			tc.Setup(t, req, server.TokenMaker)


			server.Router.ServeHTTP(res, req)

			// TODO: Test response to confirm when got the expected.
			tc.CheckResponse(t, *res)

		})
	}
}
