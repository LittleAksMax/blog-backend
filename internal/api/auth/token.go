package auth

import (
	"firebase.google.com/go/v4/auth"
	"strings"
	"time"
)

func splitBearer(bearer string) (string, bool) {
	strs := strings.Split(bearer, "Bearer ")
	if len(strs) != 2 {
		return "", false
	}

	return strs[1], true
}

func checkExpired(token *auth.Token) bool {
	return token.Expires < time.Now().Unix()
}

func CheckExists(bearer string) bool {
	return bearer != ""
}
