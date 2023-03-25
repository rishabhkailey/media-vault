package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// for now media will have 1-to-1 relationship with upload request and metadata
// todo make uploadrequest pointer so we can set the foreign key to null of media without upload request
type Media struct {
	gorm.Model
	FileName        string
	UploadRequestID string `gorm:"index"`
	UploadRequest   UploadRequest
	Metadata        MediaMetadata
	MetadataID      uint
}

type MediaModel struct {
	Db *gorm.DB
}

func NewMediaModel(db *gorm.DB) (*MediaModel, error) {
	err := db.Migrator().AutoMigrate(&Media{})
	if err != nil {
		return nil, err
	}
	return &MediaModel{
		Db: db,
	}, nil
}

func (model *MediaModel) Create(ctx context.Context, uploadRequest UploadRequest, metadata MediaMetadata) (*Media, error) {
	media := Media{
		FileName:      uuid.New().String(),
		UploadRequest: uploadRequest,
		Metadata:      metadata,
	}
	err := model.Db.WithContext(ctx).Create(&media).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return &media, nil
}

func (model *MediaModel) FindByUploadRequest(ctx context.Context, uploadRequestID string) (media Media, err error) {
	err = model.Db.First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}
