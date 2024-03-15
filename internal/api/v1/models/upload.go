package v1models

import (
	"fmt"
	"mime/multipart"
	"time"

	mediametadata "github.com/rishabhkailey/media-vault/internal/services/mediaMetadata"
)

type InitChunkUploadRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	Size      int64  `json:"size" binding:"required"`
	MediaType string `json:"media_type"`
	Date      int64  `json:"date,omitempty" binding:"required"`
}

func ValidateInitChunkUploadRequest(body InitChunkUploadRequest) (InitChunkUploadRequest, error) {
	if time.UnixMilli(body.Date).After(time.Now()) || body.Date == 0 {
		return body, fmt.Errorf("[validateInitChunkUploadRequestBody]: invalid date")
	}
	if len(body.FileName) == 0 || body.Size == 0 {
		return body, fmt.Errorf("[validateInitChunkUploadRequestBody]: invalid file metadata")
	}
	if len(body.MediaType) == 0 {
		body.MediaType = mediametadata.TYPE_UNKNOWN
	}
	return body, nil
}

type InitChunkUploadResponse struct {
	RequestID string `json:"request_id"`
	FileName  string `json:"file_name"`
}

// *int for binding to not fail for 0 value https://github.com/go-playground/validator/issues/692#issuecomment-737039536
// todo do memory usage test using big chunks
type UploadChunkRequest struct {
	Index     *int64                `form:"index" binding:"required,number,gte=0"`
	ChunkSize int64                 `form:"chunk_size" binding:"required"`
	ChunkData *multipart.FileHeader `form:"chunk_data" binding:"required"`
}

type UploadChunkResponse struct {
	RequestID string `json:"request_id" binding:"required"`
	Uploaded  int64  `json:"uploaded" binding:"required"`
}

type UploadThumbnailRequest struct {
	Size                 int64                 `form:"size" binding:"required"`
	Thumbnail            *multipart.FileHeader `form:"thumbnail" binding:"required"` //change to thumbnailData
	ThumbnailAspectRatio float32               `form:"thumbnail_aspect_ratio"`
}

type UploadThumbnailResponse struct {
	RequestID string `json:"request_id" binding:"required"`
	Uploaded  int64  `json:"uploaded" binding:"required"`
}

type FinishUploadRequest struct {
	Checksum string `json:"checksum"`
}

type FinishUploadResponse struct {
	GetMediaResponse
}
