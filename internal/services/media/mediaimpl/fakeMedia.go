package mediaimpl

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/services/media"
	"gorm.io/gorm"
)

type FakeService struct {
	ExpectedMedia      media.Model
	ExpectedMediaList  []media.GetMediaQueryResultItem
	ExpectedMediaList2 []media.Model
	ExpectedError      error
}

func NewFakeService() media.Service {
	return &FakeService{}
}

var _ media.Service = (*FakeService)(nil)

func (s *FakeService) WithTransaction(_ *gorm.DB) media.Service {
	return s
}

func (s *FakeService) Create(ctx context.Context, cmd media.CreateMediaCommand) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) DeleteOne(ctx context.Context, cmd media.DeleteOneCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteMany(ctx context.Context, cmd media.DeleteManyCommand) error {
	return s.ExpectedError
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

func (s *FakeService) GetByUserID(ctx context.Context, cmd media.GetByUserIDQuery) ([]media.Model, error) {
	return s.ExpectedMediaList2, s.ExpectedError
}

func (s *FakeService) GetByMediaIDs(ctx context.Context, query media.GetByMediaIDsQuery) ([]media.Model, error) {
	return s.ExpectedMediaList2, s.ExpectedError
}
func (s *FakeService) GetByMediaID(context.Context, media.GetByMediaIDQuery) (media.Model, error) {
	return s.ExpectedMedia, s.ExpectedError
}

func (s *FakeService) GetTypeByFileName(ctx context.Context, cmd media.GetTypeByFileNameQuery) (string, error) {
	return s.ExpectedMedia.Metadata.Type, s.ExpectedError
}
