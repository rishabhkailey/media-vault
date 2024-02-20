package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-vault/internal/api/v1/models"
)

func (server *Server) GetSpaConfig(c *gin.Context) {
	c.JSON(http.StatusOK, v1models.SpaConfig{
		OidcServerUrl:               server.config.OIDC.URL,
		OidcServerPublicClientId:    server.config.OIDC.SPA.ClientID,
		OidcServerDiscoveryEndpoint: server.config.OIDC.DiscoveryEndpoint,
	})
}
