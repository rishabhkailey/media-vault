package db

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rishabhkailey/media-service/internal/config"
)

func NewMinioConnection(config config.MinioConfig) (*minio.Client, error) {
	return minio.New(fmt.Sprintf("%s:%d", config.Host, config.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(config.User, config.Password, ""),
		Secure: config.SSL,
	})
}
