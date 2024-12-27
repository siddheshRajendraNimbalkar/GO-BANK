package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
)

const (
	AuthHeader = "Authorization"
)

// AuthMiddleWare validates the authorization token
func AuthMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthHeader)
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		authHeaderSplit := strings.Fields(authHeader)
		if len(authHeaderSplit) != 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			ctx.Abort()
			return
		}

		tokenString := authHeaderSplit[1]
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		fmt.Println(payload)

		ctx.Set("user_payload", payload)
		ctx.Next()
	}
}
