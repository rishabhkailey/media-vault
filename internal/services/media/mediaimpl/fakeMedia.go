package mediaimpl

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/services/media"
)

type FakeService struct {
	ExpectedMedia     media.Model
	ExpectedMediaList []media.GetMediaQueryResultItem
	ExpectedError     error
}

func NewFakeService() media.Service {
	return &FakeService{}
}

var _ media.Service = (*FakeService)(nil)

func (s *FakeService) Create(ctx context.Context, cmd media.CreateMediaCommand) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) GetByUploadRequestID(ctx context.Context, cmd media.GetByUploadRequestQuery) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) GetMediaWithMetadataByUploadRequestID(ctx context.Context, cmd media.GetByUploadRequestQuery) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) GetByFileName(ctx context.Context, cmd media.GetByFileNameQuery) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) GetByUserID(ctx context.Context, cmd media.GetByUserIDQuery) ([]media.GetMediaQueryResultItem, error) {
	return s.ExpectedMediaList, s.ExpectedError
}

func (s *FakeService) GetByMediaIDs(ctx context.Context, ids []uint) ([]media.GetMediaQueryResultItem, error) {
	return s.ExpectedMediaList, s.ExpectedError
}

func (s *FakeService) GetTypeByFileName(ctx context.Context, cmd media.GetTypeByFileNameQuery) (string, error) {
	return s.ExpectedMedia.Metadata.Type, s.ExpectedError
}
