package usermediabindingsimpl

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
)

type FakeService struct {
	ExpectedError             error
	ExpectedUserMediaBinding  usermediabindings.Model
	ExpectedFileBelongsToUser bool
}

var _ usermediabindings.Service = (*FakeService)(nil)

func NewFakeService() usermediabindings.Service {
	return &FakeService{}
}

func (s *FakeService) Create(ctx context.Context, query usermediabindings.CreateCommand) (usermediabindings.Model, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) GetByMediaID(ctx context.Context, query usermediabindings.GetByMediaIDQuery) (usermediabindings.Model, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) CheckFileBelongsToUser(ctx context.Context, query usermediabindings.CheckFileBelongsToUserQuery) (bool, error) {
	return s.ExpectedFileBelongsToUser, s.ExpectedError
}

func (s *FakeService) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) ([]media.Model, error) {
	return []media.Model{}, s.ExpectedError
}
