package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dassyareg/bank_app/token"
	"github.com/gin-gonic/gin"
)

const (
	AuthPayload     = "Auth_Payload"
	AuthType        = "bearer"
	AuthHeaderField = "Authorization"

	ErrNoAuth = "No Authorization Header found"
	ErrNoToken = "No Authentication token provide"
	ErrInvalidAuthType = "Invalid Auth method"
	ErrAuthFailed = "Auth failed"

)

// AuthMiddleWare returns an handlerFunc that handles authentication before
// passing the request to the intended endpoint
func AuthMiddleWare(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(gc *gin.Context) {
		authHeader := gc.GetHeader(AuthHeaderField)
		if len(authHeader) <= 0 || authHeader == "" {
			err := errors.New(ErrNoAuth)
			gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
			return
		}

		auth := strings.Fields(authHeader)

		if len(auth) < 2 {
			err := errors.New(ErrNoToken)
			gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if strings.ToLower(auth[0]) != AuthType {
			err := errors.New(ErrInvalidAuthType)
			gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token := auth[1]

		Payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			err := errors.New(ErrAuthFailed)
			gc.AbortWithStatusJSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		gc.Set(AuthPayload, Payload)

		gc.Next()

	}

}
