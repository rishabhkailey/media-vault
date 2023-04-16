package services

import (
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediametadataimpl "github.com/rishabhkailey/media-service/internal/services/mediaMetadata/mediaMetadataImpl"
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
}

func NewServices(db *gorm.DB) (*Services, error) {
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

	return &Services{
		Media:             mediaService,
		UserMediaBindings: userMediaBindingsService,
		MediaMetadata:     mediaMetadataService,
		UploadRequests:    uploadRequestsService,
	}, nil
}
