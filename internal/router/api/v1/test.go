package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi",
		"success": true,
	})
}
