package auth

import (
	"net/http"
	"strings"

	"github.com/rishabhkailey/media-service/internal/utils"
)

func GetBearerToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token, token != ""
}

func ValidateScope(userScope string, requiredScopes []string) bool {
	if len(requiredScopes) == 0 {
		return true
	}
	userScopes := strings.Split(userScope, " ")
	if len(userScopes) == 0 {
		return false
	}
	return utils.ContainsSlice(userScopes, requiredScopes)
}
