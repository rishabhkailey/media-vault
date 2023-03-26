package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/sirupsen/logrus"
)

func (server Server) UserAuthMiddleware(c *gin.Context) {
	// to make APIs behave consistent, token is required no matter what even if we are using session.
	_, ok := auth.GetBearerToken(c.Request)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userSession, err := db.GetUserSession(c)
	if err != nil {
		logrus.Errorf("[UserAuthMiddleware] get user session failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if userSession.HasUserScope && userSession.SessionExpireTime > time.Now().Unix() {
		c.Keys["userID"] = userSession.UserID
		c.Keys["userScope"] = userSession.UserScope
		c.Next()
		return
	}
	server.UserTokenAuthMiddleWare(c)
}

func (server Server) UserTokenAuthMiddleWare(c *gin.Context) {
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
	// min of 1 hour or token expire time
	// sessionExpireTime := time.Now().Add(1 * time.Hour).Unix()
	// if sessionExpireTime > tokenInfo.ExpireTime {
	// 	sessionExpireTime = tokenInfo.ExpireTime
	// }
	if err := db.SetUserSession(c, db.UserSession{
		UserID:            tokenInfo.Subject,
		UserScope:         tokenInfo.Subject,
		SessionExpireTime: tokenInfo.ExpireTime,
		HasUserScope:      true,
	}); err != nil {
		logrus.Warnf("[UserTokenAuthMiddleWare] failed to set user session: %w", err)
	}
	c.Keys["userID"] = tokenInfo.Subject
	c.Keys["userScope"] = tokenInfo.Scope
	c.Keys["expires"] = tokenInfo.ExpireTime
	c.Next()
}

func (server Server) RefreshSession(c *gin.Context) {
	expires, ok := c.Keys["expires"].(int64)
	if !ok || expires == 0 {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"expires": expires,
	})
}

// session based only, we cannot use bearer token for media files request
// on app open we will check if the token is expiring in less than 24 hours we will request new token and referesh user session to avoid interupption
// as session can expire while user is mid download/playback
func (server *Server) SessionBasedMediaAuthMiddleware(c *gin.Context) {
	userSession, err := db.GetUserSession(c)
	if err != nil {
		logrus.Errorf("[SessionBasedMediaAuthMiddleware] get user session failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !userSession.HasUserScope || userSession.SessionExpireTime < time.Now().Unix() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fileName := c.Param("fileName")
	if len(fileName) == 0 {
		logrus.Errorf("[SessionBasedMediaAuthMiddleware]: empty file name")
		c.AbortWithStatus(http.StatusBadRequest)
	}
	fileBelongsToUser, err := db.GetUserFileAccessFromSession(c, fileName, userSession.UserID)
	if err != nil {
		logrus.Errorf("[SessionBasedMediaAuthMiddleware] get user file access from session failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if fileBelongsToUser {
		c.Next()
	}
	ok, err := server.UserMediaBindings.CheckFileBelongsToUser(c.Request.Context(), userSession.UserID, fileName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.UserFileAuthMiddleware"}).Errorf("db.CheckFileBelongsToUser failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if err := db.SetUserFileAccessInSession(c, fileName, userSession.UserID); err != nil {
		logrus.Warn("[SessionBasedMediaAuthMiddleware] set user file access in session failed: %w", err)
	}
	c.Next()
}

func (server *Server) TerminateSession(c *gin.Context) {
	if err := db.DeleteUserSession(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}
