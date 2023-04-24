package v1

import (
	"net/url"

	"github.com/go-session/session/v3"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/rishabhkailey/media-service/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	Config           *config.Config
	RedisStore       *db.RedisStore
	Db               *gorm.DB
	Minio            *minio.Client
	OidcClient       auth.OidcClient
	MinioObjectCache *db.MinioObjectCache
	services.Services
}

func NewServer(config *config.Config) (*Server, error) {

	redis, err := db.NewRedisClient(config.Cache)
	if err != nil {
		return nil, err
	}

	// todo remove this
	redisStore, err := db.NewRedisStore(config.Cache)
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

	services, err := services.NewServices(DbConn, meiliSearchClient, minioClient, redis)
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

	return &Server{
		Config:           config,
		RedisStore:       redisStore,
		Db:               DbConn,
		Services:         *services,
		Minio:            minioClient,
		OidcClient:       *oidcClient,
		MinioObjectCache: db.NewMinioObjectCache(minioClient),
	}, nil
}
