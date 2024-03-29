package mediastorage

import (
	"io"
	"io/fs"
	"net/http"

	"github.com/rishabhkailey/media-vault/internal/utils"
)

type GetMediaByFileNameQuery struct {
	FileName string
}

type GetThumbnailByFileNameQuery struct {
	FileName string
}

type HttpGetRangeHandlerQuery struct {
	FileName       string
	Range          utils.Range
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HttpGetMediaHandlerQuery struct {
	FileName       string
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HttpGetThumbnailHandlerQuery struct {
	FileName       string
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type InitChunkUploadCmd struct {
	UserID    string
	RequestID string
	FileName  string
	FileSize  int64
}

type UploadChunkCmd struct {
	UserID          string
	UploadRequestID string
	Index           int64
	ChunkSize       int64
	Chunk           io.Reader
}

type FinishChunkUpload struct {
	UserID    string
	RequestID string
	CheckSum  string
}

type UploadThumbnailCmd struct {
	RequestID  string
	UserID     string
	FileName   string
	FileSize   int64
	FileReader io.Reader
}

type StoreSaveFileCmd struct {
	FileName   string
	FileSize   int64
	FileReader io.Reader
}

type File interface {
	io.ReadSeeker
	fs.File
}

type DeleteOneCommand struct {
	FileName     string
	HasThumbnail bool
}

type DeleteManyCommand struct {
	DeleteCmds []DeleteOneCommand
}
