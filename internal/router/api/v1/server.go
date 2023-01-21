package v1

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-session/session/v3"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OidcClient struct {
	provider oidc.Provider
	verfier  oidc.IDTokenVerifier
	// config   oauth2.Config // oauth config contains redirect uri and each request will have different reqiest so we can not use common config for all requests
}

func NewOidcClient(url, clientID, clientSecret string) (*OidcClient, error) {
	var err error

	oidcProvider, err := oidc.NewProvider(context.Background(), url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   url,
			"error": err,
		}).Error("failed to get oidc provider config")
		return nil, nil
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return &OidcClient{
		provider: *oidcProvider,
		verfier:  *oidcVerifier,
	}, nil
}

type Server struct {
	Config     *config.Config
	TokenStore *db.RedisStore
	Db         *gorm.DB
	Minio      *minio.Client
	OidcClient OidcClient
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

	oidcClient, err := NewOidcClient(config.AuthService.URL, config.AuthService.ID, config.AuthService.Secret)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:     config,
		TokenStore: tokenStore,
		Db:         DbConn,
		Services:   *services,
		Minio:      minioClient,
		OidcClient: *oidcClient,
	}, nil
}
