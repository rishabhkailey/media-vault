package db

import (
	dbservices "github.com/rishabhkailey/media-service/internal/db/services"
	"gorm.io/gorm"
)

type Services struct {
	Media             dbservices.MediaModel
	MediaMetadata     dbservices.MediaMetadataModel
	UserMediaBindings dbservices.UserMediaBindingModel
	UploadRequests    dbservices.UploadRequestModel
}

func NewServices(db *gorm.DB) (*Services, error) {
	// order matters, order of table creation
	uploadRequestModel, err := dbservices.NewUploadRequestModel(db)
	if err != nil {
		return nil, err
	}
	mediaModel, err := dbservices.NewMediaModel(db)
	if err != nil {
		return nil, err
	}
	mediaMetadataModel, err := dbservices.NewMediaMetadataModel(db)
	if err != nil {
		return nil, err
	}
	userMediaBindingModel, err := dbservices.NewUserMediaBinding(db)
	if err != nil {
		return nil, err
	}

	return &Services{
		Media:             *mediaModel,
		UserMediaBindings: *userMediaBindingModel,
		MediaMetadata:     *mediaMetadataModel,
		UploadRequests:    *uploadRequestModel,
	}, nil
}
