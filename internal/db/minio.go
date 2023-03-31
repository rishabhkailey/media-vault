package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewMinioConnection(config config.MinioConfig) (*minio.Client, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          256,
		MaxIdleConnsPerHost:   16,
		ResponseHeaderTimeout: time.Minute,
		IdleConnTimeout:       time.Minute,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,

		DisableCompression: true,
	}

	if config.TLS {
		tr.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		if config.CustomRootCAPath != "" {
			rootCAs, err := x509.SystemCertPool()
			if err != nil {
				logrus.Error("unable to get systme cert pool, continuing with empty cert pool")
				rootCAs = x509.NewCertPool()
			}
			data, err := os.ReadFile(config.CustomRootCAPath)
			if err == nil {
				rootCAs.AppendCertsFromPEM(data)
			}
			tr.TLSClientConfig.RootCAs = rootCAs
		}
	}
	if config.TLSSkipVerify {
		tr.TLSClientConfig.InsecureSkipVerify = true
	}
	client, err := minio.New(fmt.Sprintf("%s:%d", config.Host, config.Port), &minio.Options{
		Creds:     credentials.NewStaticV4(config.User, config.Password, ""),
		Secure:    config.TLS,
		Transport: tr,
	})
	if err != nil {
		return nil, err
	}

	return client, err
}

// todo lfu
type MinioObjectCache struct {
	cached map[string]*minio.Object
	mu     sync.Mutex
	cli    *minio.Client
}

func NewMinioObjectCache(cli *minio.Client) *MinioObjectCache {
	return &MinioObjectCache{
		cached: make(map[string]*minio.Object),
		cli:    cli,
		mu:     sync.Mutex{},
	}
}

func (cache *MinioObjectCache) Get(ctx context.Context, bucket, objectName string) (object *minio.Object, err error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cacheKey := fmt.Sprintf("%s:%s", bucket, objectName)
	object, ok := cache.cached[cacheKey]
	if !ok || object == nil {
		logrus.Warnf("cache miss bucket=%s object=%s", bucket, object)
		object, err = cache.cli.GetObject(context.Background(), bucket, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		cache.cached[cacheKey] = object
	}
	return object, nil
}
