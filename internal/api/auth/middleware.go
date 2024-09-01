package auth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	bearerTokenKey = "Authorization"
	adminClaim     = "admin"
	tokenKey       = "token"
	adminKey       = "admin"
)

func ReadToken(authClient *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token *auth.Token = nil
		var err error = nil

		bearer := ctx.GetHeader(bearerTokenKey)
		if bearer != "" { // bearer exists
			// remove 'Bearer' prefix
			tokenStr, ok := splitBearer(bearer)
			if ok { // we can split the token from the prefix (valid token format)
				token, err = authClient.VerifyIDTokenAndCheckRevoked(ctx.Request.Context(), tokenStr)
				if err != nil { // invalid token
					token = nil // just to make sure it is set to nil
				}
			}
		}

		ctx.Set(tokenKey, token)
		ctx.Next()
	}
}

func ReadAdmin(ctx *gin.Context) {
	token := ctx.MustGet(tokenKey).(*auth.Token)

	if token != nil { // token exists
		if !checkExpired(token) { // token is not expired
			admin, ok := token.Claims[adminClaim]
			adminBool, validAdminClaim := admin.(bool)

			// the admin claim must be boolean
			if !validAdminClaim {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			if ok && adminBool {
				ctx.Set(adminKey, true)
				ctx.Next()
				return
			}
		}
	}

	ctx.Set(adminKey, false)
	ctx.Next()
}

func RequiresAdmin(ctx *gin.Context) {
	admin := ctx.MustGet(adminKey).(bool)

	if !admin {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "must be an admin to access this page."})
		return
	}

	ctx.Next()
}
