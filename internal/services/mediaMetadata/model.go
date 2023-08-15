package mediametadata

import (
	"time"

	"gorm.io/gorm"
)

const (
	TABLE_NAME      = "media_metadata"
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
	return TABLE_NAME
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
