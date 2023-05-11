package mediastorageimpl

import (
	"bytes"
	"context"
	"fmt"
	"io"

	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
)

type FakeService struct {
	ExpectedError        error
	ExpectedWrittenBytes int64
	FileBytes            []byte
}

func NewFakeService() mediastorage.Service {
	return &FakeService{}
}

var _ mediastorage.Service = (*FakeService)(nil)

func (s *FakeService) HttpGetRangeHandler(ctx context.Context, query mediastorage.HttpGetRangeHandlerQuery) (int64, error) {
	if s.ExpectedError != nil {
		return s.ExpectedWrittenBytes, s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	query.ResponseWriter.Header().Add("Range", fmt.Sprintf("bytes=%d-%d/%d", query.Range.Start, query.Range.End, len(s.FileBytes)))
	return s.ExpectedWrittenBytes, s.ExpectedError
}

func (s *FakeService) HttpGetMediaHandler(ctx context.Context, query mediastorage.HttpGetMediaHandlerQuery) (int64, error) {
	if s.ExpectedError != nil {
		return s.ExpectedWrittenBytes, s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	return s.ExpectedWrittenBytes, s.ExpectedError
}

func (s *FakeService) HttpGetThumbnailHandler(ctx context.Context, query mediastorage.HttpGetThumbnailHandlerQuery) (int64, error) {
	if s.ExpectedError != nil {
		return s.ExpectedWrittenBytes, s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	return s.ExpectedWrittenBytes, s.ExpectedError
}

func (s *FakeService) InitChunkUpload(ctx context.Context, cmd mediastorage.InitChunkUploadCmd) error {
	return s.ExpectedError
}

func (s *FakeService) UploadChunk(ctx context.Context, cmd mediastorage.UploadChunkCmd) (int64, error) {
	return s.ExpectedWrittenBytes, s.ExpectedError
}

func (s *FakeService) FinishChunkUpload(ctx context.Context, cmd mediastorage.FinishChunkUpload) error {
	return s.ExpectedError
}

func (s *FakeService) ThumbnailUpload(ctx context.Context, cmd mediastorage.UploadThumbnailCmd) error {
	return s.ExpectedError
}

func (s *FakeService) GetThumbnailFileName(fileName string) string {
	return fmt.Sprintf(".thumb-%s", fileName)
}
