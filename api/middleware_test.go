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
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func AddAuthHeader(t *testing.T, username string, duration time.Duration, tokenMaker token.TokenMaker, req *http.Request, AuthorizationType string) {
	Token, err :=  tokenMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	AuthValue := fmt.Sprintf("%s %s", AuthorizationType, Token)
	req.Header.Set(AuthHeaderField, AuthValue)
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
				AddAuthHeader(t, utils.RandomName(), time.Minute, tokenMaker, req, AuthType)
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusOK)
			},
		},
		{
			Name: "Expired token",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, utils.RandomName(), -time.Minute, tokenMaker, req, AuthType)
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusUnauthorized)
				expectedErr, err := json.Marshal(gin.H{"error": ErrAuthFailed})
				require.NoError(t, err)
				require.Equal(t, expectedErr, res.Body.Bytes())
			},
		},
		{
			Name: "Wrong auth type",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, utils.RandomName(), time.Minute, tokenMaker, req, "AuthType")
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
				expectedErr, err := json.Marshal(gin.H{"error": ErrInvalidAuthType})
				require.NoError(t, err)
				require.Equal(t, expectedErr, res.Body.Bytes())
			},
		},
		{
			Name: "No auth type",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				AddAuthHeader(t, utils.RandomName(), time.Minute, tokenMaker, req, "")
	
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
				expectedErr, err := json.Marshal(gin.H{"error": ErrNoToken})
				require.NoError(t, err)
				require.Equal(t, expectedErr, res.Body.Bytes())
			},
		},
		{
			Name: "Wrong Token",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
				CreateToken, err := token.NewPaseToken(utils.RdmString(32))
				require.NoError(t,err)

				Token, err := CreateToken.GenerateToken("Memphis", time.Hour)
				require.NoError(t, err)
	
				AuthValue := fmt.Sprintf("%s %s", AuthType, Token)
	
				req.Header.Set(AuthHeaderField, AuthValue)
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusUnauthorized)
				expectedErr, err := json.Marshal(gin.H{"error": ErrAuthFailed})
				require.NoError(t, err)
				require.Equal(t, expectedErr, res.Body.Bytes())
			},
		},
		{
			Name: "No Token",
	
			Setup: func(t *testing.T, req *http.Request, tokenMaker token.TokenMaker) {
	
			},
			CheckResponse: func(t *testing.T, res httptest.ResponseRecorder) {
				require.Equal(t, res.Code, http.StatusBadRequest)
				expectedErr, err := json.Marshal(gin.H{"error": ErrNoAuth})
				require.NoError(t, err)
				require.Equal(t, expectedErr, res.Body.Bytes())
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

			tc.Setup(t, req, server.TokenMaker)


			server.Router.ServeHTTP(res, req)

			tc.CheckResponse(t, *res)

		})
	}
}
