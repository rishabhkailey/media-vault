package v1

import (
	"net/url"

	"github.com/go-session/session/v3"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/rishabhkailey/media-service/internal/services"
)

// todo interface so we can test it

// rename to something else? Server doesn't seems right
// api ?
type Server struct {
	services.Services
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
	session.InitManager(
		session.SetStore(db.NewRedisSessionStore(config.Cache)),
	)

	meiliSearchClient, err := db.NewMeiliSearchClient(config.MeiliSearch)
	if err != nil {
		return nil, err
	}

	minioClient, err := db.NewMinioConnection(config.MinioConfig)
	if err != nil {
		return nil, err
	}

	redirectURI, err := url.JoinPath(config.Server.BaseURL, "/v1/authorize")
	if err != nil {
		return nil, err
	}

	oidcClient, err := auth.NewOidcClient(config.AuthService.URL, config.AuthService.ClientID, config.AuthService.ClientSecret, redirectURI)
	if err != nil {
		return nil, err
	}

	services, err := services.NewServices(DbConn, meiliSearchClient, minioClient, redisClient, oidcClient)
	if err != nil {
		return nil, err
	}

	return &Server{
		Services: *services,
	}, nil
}
