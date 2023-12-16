package mediastorage

import (
	"context"
)

type Service interface {
	// GetMediaByFileName(context.Context, GetMediaByFileNameQuery) (File, error)
	// GetThumbnailByFileName(context.Context, GetThumbnailByFileNameQuery) (File, error)
	HttpGetRangeHandler(context.Context, HttpGetRangeHandlerQuery) error
	HttpGetMediaHandler(context.Context, HttpGetMediaHandlerQuery) error
	HttpGetThumbnailHandler(context.Context, HttpGetThumbnailHandlerQuery) error
	InitChunkUpload(context.Context, InitChunkUploadCmd) error
	UploadChunk(context.Context, UploadChunkCmd) (int64, error)
	FinishChunkUpload(context.Context, FinishChunkUpload) error
	ThumbnailUpload(context.Context, UploadThumbnailCmd) error
	GetThumbnailFileName(string) string
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) (failedFileNames []string, errs []error)
}
