package v1

import (
	"github.com/go-session/session/v3"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/db"
	"gorm.io/gorm"
)

type Server struct {
	Config     *config.Config
	TokenStore *db.RedisStore
	Db         *gorm.DB
	Minio      *minio.Client
	db.Services
}

func NewServer(config *config.Config) (*Server, error) {

	tokenStore, err := db.NewRedisTokenStore(config.Cache)
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

	services, err := db.NewServices(DbConn)
	if err != nil {
		return nil, err
	}

	minioClient, err := db.NewMinioConnection(config.MinioConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:     config,
		TokenStore: tokenStore,
		Db:         DbConn,
		Services:   *services,
		Minio:      minioClient,
	}, nil
}
