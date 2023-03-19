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

const VIDEO_MP4 MediaType = "video/webm"
const VIDEO_WEBM MediaType = "video/webm"
const IMAGE_PNG MediaType = "image/png"
const IMAGE_JPEG MediaType = "image/jpeg"
const UNKNOWN MediaType = "unknown"

var ValidMediaTypes = []MediaType{VIDEO_MP4, VIDEO_WEBM, IMAGE_JPEG, IMAGE_PNG, UNKNOWN}

type Metadata struct {
	Name string
	Date time.Time
	Type uint64
	Size uint64
}

type MediaMetadata struct {
	gorm.Model
	Metadata
	MediaID uint `gorm:"index"`
	Media   Media
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

func (model *MediaMetadataModel) Create(ctx context.Context, media Media, metadata Metadata) (*MediaMetadata, error) {
	mediaMetadata := MediaMetadata{
		Metadata: metadata,
		Media:    media,
	}
	err := model.Db.WithContext(ctx).Create(&mediaMetadata).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return &mediaMetadata, nil
}
