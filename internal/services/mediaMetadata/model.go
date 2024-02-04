package mediametadata

import (
	"time"

	"github.com/rishabhkailey/media-service/internal/constants"
	"gorm.io/gorm"
)

const (
	TYPE_UNKNOWN    = "unknown"
	TYPE_VIDEO_MP4  = "video/mp4"
	TYPE_VIDEO_WEBM = "video/webm"
	TYPE_IMAGE_PNG  = "image/png"
	TYPE_IMAGE_JPEG = "image/jpeg"
)

type Metadata struct {
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Type                 string    `json:"type"`
	Size                 uint64    `json:"size"`
	Thumbnail            bool      `gorm:"default:false" json:"thumbnail"`
	ThumbnailAspectRatio float32   `json:"thumbnail_aspect_ratio"`
}

type Model struct {
	gorm.Model
	Metadata
}

func (Model) TableName() string {
	return constants.MEDIA_METADATA_TABLE
}

type CreateCommand struct {
	Metadata
}

type UpdateThumbnailCommand struct {
	ID                   uint
	Thumbnail            bool
	ThumbnailAspectRatio float32
}

type DeleteOneCommand struct {
	ID uint
}

type DeleteManyCommand struct {
	IDs []uint
}
