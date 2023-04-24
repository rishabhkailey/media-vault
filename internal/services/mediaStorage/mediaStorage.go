package mediastorage

import (
	"context"
)

type Service interface {
	GetMediaByFileName(context.Context, GetMediaByFileNameQuery) (File, error)
	GetThumbnailByFileName(context.Context, GetThumbnailByFileNameQuery) (File, error)
	HttpGetRangeHandler(context.Context, WriteRangeByFileNameQuery) (int64, error)
	HttpGetMediaHandler(context.Context, HttpGetMediaHandlerQuery) (int64, error)
	HttpGetThumbnailHandler(context.Context, HttpGetThumbnailHandlerQuery) (int64, error)
	InitChunkUpload(context.Context, InitChunkUploadCmd) error
	UploadChunk(context.Context, UploadChunkCmd) (int64, error)
	FinishChunkUpload(context.Context, FinishChunkUpload) error
	ThumbnailUpload(context.Context, UploadThumbnailCmd) error
	GetThumbnailFileName(string) string
}
