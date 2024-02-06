package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	authservice "github.com/rishabhkailey/media-service/internal/services/authService"
)

func (server *Server) UserAuthMiddleware(c *gin.Context) {
	userID, err := server.AuthService.ValidateUserAccess(authservice.ValidateUserAccessQuery{
		SessionStoreQuery: authservice.SessionStoreQuery{
			Ctx:            c.Request.Context(),
			ResponseWriter: c.Writer,
			Request:        c.Request,
		},
	}, []string{authservice.UserScope})
	if errors.Is(err, authservice.ErrUnauthorized) {
		c.Abort()
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	if errors.Is(err, authservice.ErrForbidden) {
		c.Abort()
		c.Error(internalErrors.NewForbiddenError(err))
		return
	}
	if err != nil {
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(err))
		return
	}
	c.Set("user_id", userID)
	c.Next()
}

// refresh session endpoint called on application startup so user can watch media without interuption
func (server *Server) RefreshSession(c *gin.Context) {
	expires, err := server.AuthService.RefreshSession(authservice.RefreshSessionQuery{
		SessionStoreQuery: authservice.SessionStoreQuery{
			Ctx:            c.Request.Context(),
			ResponseWriter: c.Writer,
			Request:        c.Request,
		},
	})
	// todo instead of this duplicate code, we can use errors from internalErrors package and no need to check error type here
	if errors.Is(err, authservice.ErrUnauthorized) {
		c.Abort()
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	if errors.Is(err, authservice.ErrForbidden) {
		c.Abort()
		c.Error(internalErrors.NewForbiddenError(err))
		return
	}
	if err != nil {
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[RefreshSession] get session expire time failed: %w", err),
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
func (server *Server) SessionBasedMediaFileAuthMiddleware(c *gin.Context) {
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
	err := server.AuthService.ValidateUserMediaAccess(authservice.ValidateUserMediaAccessQuery{
		SessionStoreQuery: authservice.SessionStoreQuery{
			Ctx:            c.Request.Context(),
			ResponseWriter: c.Writer,
			Request:        c.Request,
		},
		FileName: fileName,
	})
	if errors.Is(err, authservice.ErrUnauthorized) {
		c.Abort()
		c.Error(internalErrors.ErrUnauthorized)
		return
	}
	if errors.Is(err, authservice.ErrForbidden) {
		c.Abort()
		c.Error(internalErrors.NewForbiddenError(err))
		return
	}
	if errors.Is(err, authservice.ErrForbidden) {
		c.Abort()
		c.Error(internalErrors.NewForbiddenError(err))
		return
	}
	if err != nil {
		c.Abort()
		c.Error(internalErrors.NewInternalServerError(err))
		return
	}
	c.Next()
}

func (server *Server) TerminateSession(c *gin.Context) {
	err := server.AuthService.TerminateSession(authservice.TerminateSessionCmd{
		SessionStoreQuery: authservice.SessionStoreQuery{
			Ctx:            c.Request.Context(),
			ResponseWriter: c.Writer,
			Request:        c.Request,
		},
	})
	if err != nil {
		c.Abort()
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[TerminateSession] Failed to delete user session: %w", err),
			),
		)
	}
	c.Status(http.StatusOK)
}
