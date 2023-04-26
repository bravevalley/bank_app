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
)

// AuthMiddleWare returns an handlerFunc that handles authentication before
// passing the request to the intended endpoint
func AuthMiddleWare(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(gc *gin.Context) {
		authHeader := gc.GetHeader(AuthHeaderField)
		err := errors.New("No Authentication token provide")
		if len(authHeader) <= 0 || authHeader == "" {
			gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
			return
		}

		auth := strings.Fields(authHeader)

		if len(auth) < 2 {
			err := errors.New("No Authentication token provide")
			if len(authHeader) <= 0 || authHeader == "" {
				gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
				return
			}
		}

		if strings.ToLower(auth[0]) != AuthType {
			err = errors.New("Invalid Auth method")
			gc.AbortWithStatusJSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token := auth[1]

		Payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			gc.AbortWithStatusJSON(http.StatusBadRequest, errors.New("Auth failed"))
			return
		}

		gc.Set(AuthPayload, Payload)

		gc.Next()

	}

}
