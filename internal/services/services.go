package services

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediametadataimpl "github.com/rishabhkailey/media-service/internal/services/mediaMetadata/mediaMetadataImpl"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	mediasearchimpl "github.com/rishabhkailey/media-service/internal/services/mediaSearch/mediaSearchimpl"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"github.com/rishabhkailey/media-service/internal/services/uploadRequests/uploadrequestsimpl"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	usermediabindingsimpl "github.com/rishabhkailey/media-service/internal/services/userMediaBindings/userMediaBindingsimpl"
	"gorm.io/gorm"
)

type Services struct {
	Media             media.Service
	MediaMetadata     mediametadata.Service
	UserMediaBindings usermediabindings.Service
	UploadRequests    uploadrequests.Service
	MediaSearch       mediasearch.Service
}

func NewServices(db *gorm.DB, ms *meilisearch.Client) (*Services, error) {
	// order matters, order of table creation
	uploadRequestsService, err := uploadrequestsimpl.NewService(db)
	if err != nil {
		return nil, err
	}
	mediaService, err := mediaimpl.NewService(db)
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
	return &Services{
		Media:             mediaService,
		UserMediaBindings: userMediaBindingsService,
		MediaMetadata:     mediaMetadataService,
		UploadRequests:    uploadRequestsService,
		MediaSearch:       mediaSearchService,
	}, nil
}
