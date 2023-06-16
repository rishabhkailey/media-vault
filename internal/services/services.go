package services

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/meilisearch/meilisearch-go"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/auth"
	"github.com/rishabhkailey/media-service/internal/services/album"
	"github.com/rishabhkailey/media-service/internal/services/album/albumimpl"
	authservice "github.com/rishabhkailey/media-service/internal/services/authService"
	authserviceimpl "github.com/rishabhkailey/media-service/internal/services/authService/authServiceImpl"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediametadataimpl "github.com/rishabhkailey/media-service/internal/services/mediaMetadata/mediaMetadataImpl"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	mediasearchimpl "github.com/rishabhkailey/media-service/internal/services/mediaSearch/mediaSearchimpl"
	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
	"github.com/rishabhkailey/media-service/internal/services/mediaStorage/mediastorageimpl"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"github.com/rishabhkailey/media-service/internal/services/uploadRequests/uploadrequestsimpl"
	userinfo "github.com/rishabhkailey/media-service/internal/services/userInfo"
	userinfoimpl "github.com/rishabhkailey/media-service/internal/services/userInfo/impl"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	usermediabindingsimpl "github.com/rishabhkailey/media-service/internal/services/userMediaBindings/userMediaBindingsimpl"
	"github.com/rishabhkailey/media-service/internal/store"
	"gorm.io/gorm"
)

type Services struct {
	db                *gorm.DB
	Media             media.Service
	MediaMetadata     mediametadata.Service
	UserMediaBindings usermediabindings.Service
	UploadRequests    uploadrequests.Service
	MediaSearch       mediasearch.Service
	MediaStorage      mediastorage.Service
	AuthService       authservice.Service
	AlbumService      album.Service
	UserInfoService   userinfo.Service
}

func NewServices(
	db *gorm.DB,
	ms *meilisearch.Client,
	minio *minio.Client,
	redis *redis.Client,
	oidcClient *auth.OidcClient,
) (*Services, error) {
	// order matters, order of table creation
	uploadRequestsService, err := uploadrequestsimpl.NewService(db)
	if err != nil {
		return nil, err
	}
	mediaService, err := mediaimpl.NewService(db, redis)
	if err != nil {
		return nil, err
	}
	mediaMetadataService, err := mediametadataimpl.NewService(db)
	if err != nil {
		return nil, err
	}
	userMediaBindingsService, err := usermediabindingsimpl.NewService(db)
	if err != nil {
		return nil, err
	}
	mediaSearchService, err := mediasearchimpl.NewService(ms)
	if err != nil {
		return nil, err
	}
	authService, err := authserviceimpl.NewService(*oidcClient, userMediaBindingsService, time.Hour*12)
	if err != nil {
		return nil, err
	}
	store, err := store.NewStore(db, redis)
	if err != nil {
		return nil, err
	}
	albumService, err := albumimpl.NewService(*store)
	if err != nil {
		return nil, err
	}
	userInfoService, err := userinfoimpl.NewService(*store)
	if err != nil {
		return nil, err
	}

	// todo move bucket name to config
	mediaStorageService, err := mediastorageimpl.NewMinioService(minio, "test", uploadRequestsService)
	if err != nil {
		return nil, err
	}
	return &Services{
		db:                db,
		Media:             mediaService,
		UserMediaBindings: userMediaBindingsService,
		MediaMetadata:     mediaMetadataService,
		UploadRequests:    uploadRequestsService,
		MediaSearch:       mediaSearchService,
		MediaStorage:      mediaStorageService,
		AuthService:       authService,
		AlbumService:      albumService,
		UserInfoService:   userInfoService,
	}, nil
}

// todo for mock? return nil?
func (s *Services) CreateTransaction() *gorm.DB {
	return s.db.Begin()
}
