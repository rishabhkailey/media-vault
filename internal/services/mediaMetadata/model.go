package mediametadata

import (
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

const (
	TYPE_UNKNOWN    = "unknown"
	TYPE_VIDEO_MP4  = "video/mp4"
	TYPE_VIDEO_WEBM = "video/webm"
	TYPE_IMAGE_PNG  = "image/png"
	TYPE_IMAGE_JPEG = "image/jpeg"
)

type CreateCommand struct {
	storemodels.MediaMetadata
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
