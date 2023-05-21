package mediametadataimpl

import (
	"context"

	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	"gorm.io/gorm"
)

type FakeService struct {
	ExpectedMediaMetadata mediametadata.Model
	ExpectedError         error
}

func NewFakeService() mediametadata.Service {
	return &FakeService{}
}

var _ mediametadata.Service = (*FakeService)(nil)

func (s *FakeService) WithTransaction(_ *gorm.DB) mediametadata.Service {
	return s
}

func (s *FakeService) Create(ctx context.Context, cmd mediametadata.CreateCommand) (mediametadata.Model, error) {
	return s.ExpectedMediaMetadata, s.ExpectedError
}

func (s *FakeService) DeleteOne(_ context.Context, _ mediametadata.DeleteOneCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteMany(_ context.Context, _ mediametadata.DeleteManyCommand) error {
	return s.ExpectedError
}

func (s *FakeService) UpdateThumbnail(ctx context.Context, cmd mediametadata.UpdateThumbnailCommand) error {
	return s.ExpectedError
}
