package mediastorageimpl

import (
	"bytes"
	"context"
	"fmt"
	"io"

	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
)

type FakeService struct {
	ExpectedError                     error
	ExpectedWrittenBytes              int64
	ExpectedDeleteManyFailedFileNames []string
	ExpectedDeleteManyErrs            []error
	FileBytes                         []byte
}

func NewFakeService() mediastorage.Service {
	return &FakeService{}
}

var _ mediastorage.Service = (*FakeService)(nil)

func (s *FakeService) HttpGetRangeHandler(ctx context.Context, query mediastorage.HttpGetRangeHandlerQuery) error {
	if s.ExpectedError != nil {
		return s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	query.ResponseWriter.Header().Add("Range", fmt.Sprintf("bytes=%d-%d/%d", query.Range.Start, query.Range.End, len(s.FileBytes)))
	return s.ExpectedError
}

func (s *FakeService) HttpGetMediaHandler(ctx context.Context, query mediastorage.HttpGetMediaHandlerQuery) error {
	if s.ExpectedError != nil {
		return s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	return s.ExpectedError
}

func (s *FakeService) HttpGetThumbnailHandler(ctx context.Context, query mediastorage.HttpGetThumbnailHandlerQuery) error {
	if s.ExpectedError != nil {
		return s.ExpectedError
	}
	reader := bytes.NewReader(s.FileBytes)
	io.CopyN(query.ResponseWriter, reader, s.ExpectedWrittenBytes)
	return s.ExpectedError
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

func (s *FakeService) DeleteOne(ctx context.Context, cmd mediastorage.DeleteOneCommand) error {
	return s.ExpectedError
}
func (s *FakeService) DeleteMany(ctx context.Context, cmd mediastorage.DeleteManyCommand) ([]string, []error) {
	return s.ExpectedDeleteManyFailedFileNames, s.ExpectedDeleteManyErrs
}

func (s *FakeService) GetThumbnailFileName(fileName string) string {
	return fmt.Sprintf(".thumb-%s", fileName)
}
