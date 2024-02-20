package v1

import (
	"github.com/go-session/session/v3"
	"github.com/rishabhkailey/media-vault/internal/auth"
	"github.com/rishabhkailey/media-vault/internal/config"
	"github.com/rishabhkailey/media-vault/internal/db"
	"github.com/rishabhkailey/media-vault/internal/services"
)

// todo interface so we can test it

// rename to something else? Server doesn't seems right
// api ?
type Server struct {
	services.Services
	config *config.Config
}

func NewServer(config *config.Config) (*Server, error) {

	redisClient, err := db.NewRedisClient(config.Cache)
	if err != nil {
		return nil, err
	}

	DbConn, err := db.NewGoOrmConnection(config.Database)
	if err != nil {
		return nil, err
	}

	// persistent session store
	// session.
	session.InitManager(
		session.SetStore(db.NewRedisSessionStore(config.Cache)),
		session.SetCookieName("media_service"),
	)

	meiliSearchClient, err := db.NewMeiliSearchClient(config.MeiliSearch)
	if err != nil {
		return nil, err
	}

	minioClient, err := db.NewMinioConnection(config.MinioConfig)
	if err != nil {
		return nil, err
	}

	oidcClient, err := auth.NewOidcClient(config.OIDC.URL, config.OIDC.DiscoveryEndpoint, config.OIDC.MediaVault.ClientID, config.OIDC.MediaVault.ClientSecret)
	if err != nil {
		return nil, err
	}

	services, err := services.NewServices(DbConn, meiliSearchClient, minioClient, config.MinioConfig.Bucket, redisClient, oidcClient)
	if err != nil {
		return nil, err
	}

	return &Server{
		Services: *services,
		config:   config,
	}, nil
}
