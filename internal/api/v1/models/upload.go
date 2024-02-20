package v1models

import (
	"fmt"
	"mime/multipart"
	"time"

	mediametadata "github.com/rishabhkailey/media-vault/internal/services/mediaMetadata"
)

type InitChunkUploadRequest struct {
	FileName  string `json:"fileName" binding:"required"`
	Size      int64  `json:"size" binding:"required"`
	MediaType string `json:"mediaType"`
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
	RequestID string `json:"requestID"`
	FileName  string `json:"file_name"`
}

// *int for binding to not fail for 0 value https://github.com/go-playground/validator/issues/692#issuecomment-737039536
// todo do memory usage test using big chunks
type UploadChunkRequest struct {
	RequestID string                `form:"requestID" binding:"required"`
	Index     *int64                `form:"index" binding:"required,number,gte=0"`
	ChunkSize int64                 `form:"chunkSize" binding:"required"`
	ChunkData *multipart.FileHeader `form:"chunkData" binding:"required"`
}

type UploadChunkResponse struct {
	RequestID string `json:"requestID" binding:"required"`
	Uploaded  int64  `json:"uploaded" binding:"required"`
}

type UploadThumbnailRequest struct {
	RequestID            string                `form:"requestID" binding:"required"`
	Size                 int64                 `form:"size" binding:"required"`
	Thumbnail            *multipart.FileHeader `form:"thumbnail" binding:"required"` //change to thumbnailData
	ThumbnailAspectRatio float32               `form:"thumbnail_aspect_ratio"`
}

type UploadThumbnailResponse struct {
	RequestID string `json:"requestID" binding:"required"`
	Uploaded  int64  `json:"uploaded" binding:"required"`
}

type FinishUploadRequest struct {
	RequestID string `json:"requestID" binding:"required"`
	Checksum  string `json:"checksum"`
}

type FinishUploadResponse struct {
	GetMediaResponse
}
