package storemodels

import (
	"time"

	"github.com/rishabhkailey/media-vault/internal/constants"
	"gorm.io/gorm"
)

type MediaMetadata struct {
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Type                 string    `json:"type"`
	Size                 uint64    `json:"size"`
	Thumbnail            bool      `gorm:"default:false" json:"thumbnail"`
	ThumbnailAspectRatio float32   `json:"thumbnail_aspect_ratio"`
}

type MediaMetadataModel struct {
	gorm.Model
	MediaMetadata
}

func (MediaMetadataModel) TableName() string {
	return constants.MEDIA_METADATA_TABLE
}
