package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/sirupsen/logrus"
)

func (server Server) UserAuthMiddleware(c *gin.Context) {
	// token is required no matter what even if we are using session to make APIs behave consistent
	_, ok := auth.GetBearerToken(c.Request)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.UserAuthMiddleware"}).Errorf("session start failed")
		c.Status(http.StatusInternalServerError)
		return
	}
	// check session if user is already authenticated
	var hasUserScope bool
	if value, ok := store.Get("hasUserScope"); ok {
		hasUserScope, _ = value.(bool)
	}
	var sessionExpireTime int64
	if value, ok := store.Get("sessionExpireTime"); ok {
		// ignoring errors here because on error the default sessionExpirteTime will be 0 which is not valid
		svalue, _ := value.(string)
		if len(svalue) != 0 {
			sessionExpireTime, _ = strconv.ParseInt(svalue, 10, 64)
		}
	}
	var userID string
	if value, ok := store.Get("userID"); ok {
		userID, _ = value.(string)
	}
	var userScope string
	if value, ok := store.Get("userScope"); ok {
		userScope, _ = value.(string)
	}
	if hasUserScope && sessionExpireTime > time.Now().Unix() {
		c.Keys["userID"] = userID
		c.Keys["userScope"] = userScope
		c.Next()
		return
	}
	server.UserTokenAuthMiddleWare(c)
}

func (server Server) UserTokenAuthMiddleWare(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.UserAuthMiddleware"}).Errorf("session start failed")
		c.Status(http.StatusInternalServerError)
		return
	}
	token, ok := auth.GetBearerToken(c.Request)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenInfo, err := server.OidcClient.IntrospectToken(token)
	if errors.Is(err, auth.ErrUnauthorized) || !tokenInfo.ValidateScope(auth.SCOPE_USER) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if len(tokenInfo.Subject) == 0 {
		logrus.Errorf("token info doesn't contain user info")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	store.Set("userID", tokenInfo.Subject)
	store.Set("userScope", tokenInfo.Scope)
	store.Set("hasUserScope", true)
	// min of 1 hour or token expire time
	sessionExpireTime := time.Now().Add(1 * time.Hour).Unix()
	if sessionExpireTime > tokenInfo.ExpireTime {
		sessionExpireTime = tokenInfo.ExpireTime
	}
	// we will store the unix time as string, as the values are stored in redis we looses the actual type for the value
	// int64 is converted to float64 when jsonunmarshall is called check these 2 issues for better understanding https://github.com/json-iterator/go/issues/351 and https://github.com/json-iterator/go/issues/145
	store.Set("sessionExpireTime", strconv.FormatInt(sessionExpireTime, 10))
	store.Save()

	c.Keys["userID"] = tokenInfo.Subject
	c.Keys["userScope"] = tokenInfo.Scope
	c.Next()
}

func (server Server) TestProtectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi..",
	})
}
