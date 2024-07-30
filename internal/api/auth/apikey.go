package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const apiKeyHeaderKey = "X-Api-Key"

func RequiresAPIKey(apiKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if apiKey != ctx.GetHeader(apiKeyHeaderKey) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required and must be valid."})
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}

}
