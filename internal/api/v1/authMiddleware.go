package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/rishabhkailey/media-service/internal/db"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/sirupsen/logrus"
)

func (server Server) UserAuthMiddleware(c *gin.Context) {
	// to make APIs behave consistent, token is required no matter what even if we are using session.
	_, ok := auth.GetBearerToken(c.Request)
	if !ok {
		c.Abort()
		c.Error(internalErrors.ErrMissingBearerToken)
		return
	}
	userSession, err := db.GetUserSession(c)
	if err != nil {
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[UserAuthMiddleware] get user session failed: %w", err),
		))
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
		c.Abort()
		c.Error(internalErrors.ErrMissingBearerToken)
		return
	}
	tokenInfo, err := server.OidcClient.IntrospectToken(token)
	if errors.Is(err, auth.ErrUnauthorized) || !tokenInfo.ValidateScope(auth.SCOPE_USER) {
		c.Abort()
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	if len(tokenInfo.Subject) == 0 {
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("token info returned by auth server didn't include user info"),
		))
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
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[RefreshSession] expected expire key missing in context"),
		))
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
		c.Abort()
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[SessionBasedMediaAuthMiddleware] get user session failed: %w", err),
			),
		)
		return
	}
	if !userSession.HasUserScope || userSession.SessionExpireTime < time.Now().Unix() {
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	fileName := c.Param("fileName")
	// do we need to check this?
	if len(fileName) == 0 {
		c.Abort()
		c.Error(
			internalErrors.New(
				http.StatusBadRequest,
				"[SessionBasedMediaAuthMiddleware]: empty file name",
				"file name missing",
			),
		)
		return
	}
	fileBelongsToUser, err := db.GetUserFileAccessFromSession(c, fileName, userSession.UserID)
	if err != nil {
		c.Abort()
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[SessionBasedMediaAuthMiddleware] get user file access from session failed: %w", err),
			),
		)
		return
	}
	if fileBelongsToUser {
		c.Next()
		return
	}
	ok, err := server.UserMediaBindings.CheckFileBelongsToUser(c.Request.Context(), usermediabindings.CheckFileBelongsToUserQuery{
		UserID:   userSession.UserID,
		FileName: fileName,
	})
	if err != nil {
		c.Abort()
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[SessionBasedMediaAuthMiddleware] CheckFileBelongsToUser failed: %w", err),
			),
		)
		return
	}
	if !ok {
		c.Abort()
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	if err := db.SetUserFileAccessInSession(c, fileName, userSession.UserID); err != nil {
		logrus.Warn("[SessionBasedMediaAuthMiddleware] set user file access in session failed: %w", err)
	}
	c.Next()
}

func (server *Server) TerminateSession(c *gin.Context) {
	if err := db.DeleteUserSession(c); err != nil {
		c.Abort()
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[TerminateSession] Failed to delete user session: %w", err),
			),
		)
	}
	c.Status(http.StatusOK)
}
