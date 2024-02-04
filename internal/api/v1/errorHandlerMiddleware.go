package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/sirupsen/logrus"
)

// todo not so redable?
func ErrorHandler(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}
	switch e := err.Err.(type) {
	case internalErrors.CustomError:
		if e.Status == 0 {
			e.Status = http.StatusInternalServerError
		}
		if len(e.PublicMessage) == 0 {
			e.PublicMessage = "Internal server error"
		}
		c.AbortWithStatusJSON(e.Status, gin.H{
			"error": e.PublicMessage,
		})
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"url":    c.Request.URL,
		}).Error(e.Err)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
}
