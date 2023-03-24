package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Containers
type MediaType string

// todo more or custom media type from client or dynamic from file extension
const VIDEO_MP4 MediaType = "video/mp4"
const VIDEO_WEBM MediaType = "video/webm"
const IMAGE_PNG MediaType = "image/png"
const IMAGE_JPEG MediaType = "image/jpeg"
const UNKNOWN MediaType = "unknown"

var ValidMediaTypes = []MediaType{VIDEO_MP4, VIDEO_WEBM, IMAGE_JPEG, IMAGE_PNG, UNKNOWN}

// todo thumbnail bool
type Metadata struct {
	Name      string
	Date      time.Time
	Type      string
	Size      uint64
	Thumbnail bool `gorm:"default:false"`
}

type MediaMetadata struct {
	gorm.Model
	Metadata
}

type MediaMetadataModel struct {
	Db *gorm.DB
}

func NewMediaMetadataModel(db *gorm.DB) (*MediaMetadataModel, error) {
	err := db.Migrator().AutoMigrate(&MediaMetadata{})
	if err != nil {
		return nil, err
	}
	return &MediaMetadataModel{
		Db: db,
	}, nil
}

func (model *MediaMetadataModel) Create(ctx context.Context, metadata Metadata) (*MediaMetadata, error) {
	mediaMetadata := MediaMetadata{
		Metadata: metadata,
	}
	err := model.Db.WithContext(ctx).Create(&mediaMetadata).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return &mediaMetadata, nil
}

func (model *MediaMetadataModel) UpdateThumbnail(ctx context.Context, metadataID uint, thumbnail bool) error {
	return model.Db.Model(&MediaMetadata{
		Model: gorm.Model{
			ID: metadataID,
		},
	}).Update("thumbnail", thumbnail).Error
}
