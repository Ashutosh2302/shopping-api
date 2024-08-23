package middlewares

import (
	"fmt"
	"net/http"
	"shopping_api/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized!"})
		return
	}

	claims, err := utils.VerifyAccessToken(token)

	if err != nil {
		fmt.Print("Error validating token: ", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token!"})
		return
	}

	ctx.Set("userId", claims.UserId)
	ctx.Next()
}
