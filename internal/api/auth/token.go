package auth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	bearerTokenKey = "Authorization"
	adminClaim     = "admin"
)

func splitBearer(bearer string) (string, bool) {
	strs := strings.Split(bearer, "Bearer ")
	if len(strs) != 2 {
		return "", false
	}

	return strs[1], true
}

func RequiresToken(authClient *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader(bearerTokenKey)
		if bearer == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing 'Authorization' header with bearer token"})
			return
		}

		// remove 'Bearer' prefix
		tokenStr, ok := splitBearer(bearer)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		token, err := authClient.VerifyIDTokenAndCheckRevoked(ctx.Request.Context(), tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("token", token)
		ctx.Next()
	}
}

func checkExpired(token *auth.Token) bool {
	return token.Expires < time.Now().Unix()
}

func RequiresAdmin(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	// check if token is expired
	if checkExpired(token) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	// check if we have the right claims
	admin, ok := token.Claims[adminClaim]

	// admin claim doesn't exist or isn't true
	if !ok || !(admin.(bool)) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	ctx.Next()
}
