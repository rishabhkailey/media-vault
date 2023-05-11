package mediasearchimpl

import (
	"context"

	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
)

type FakeService struct {
	ExpectedTaskID       int64
	ExpectedMediaList    []mediasearch.Model
	ExpectedMediaListIDs []uint
	ExpectedError        error
}

func NewFakeService() mediasearch.Service {
	return &FakeService{}
}

var _ mediasearch.Service = (*FakeService)(nil)

func (s *FakeService) CreateOne(ctx context.Context, cmd mediasearch.Model) (int64, error) {
	return s.ExpectedTaskID, s.ExpectedError
}

func (s *FakeService) CreateMany(ctx context.Context, cmd []mediasearch.Model) (int64, error) {
	return s.ExpectedTaskID, s.ExpectedError
}

func (s *FakeService) Search(ctx context.Context, query mediasearch.MediaSearchQuery) ([]mediasearch.Model, error) {
	return s.ExpectedMediaList, s.ExpectedError
}

func (s *FakeService) SearchGetMediaIDsOnly(ctx context.Context, query mediasearch.MediaSearchQuery) ([]uint, error) {
	return s.ExpectedMediaListIDs, s.ExpectedError
}
