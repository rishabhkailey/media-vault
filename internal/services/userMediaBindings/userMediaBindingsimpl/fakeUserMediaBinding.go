package usermediabindingsimpl

import (
	"context"

	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	mediaStore "github.com/rishabhkailey/media-service/internal/store/media"
	"gorm.io/gorm"
)

type FakeService struct {
	ExpectedError              error
	ExpectedUserMediaBinding   usermediabindings.Model
	ExpectedFileBelongsToUser  bool
	ExpectedMediaBelongsToUser bool
}

var _ usermediabindings.Service = (*FakeService)(nil)

func NewFakeService() usermediabindings.Service {
	return &FakeService{}
}

func (s *FakeService) WithTransaction(_ *gorm.DB) usermediabindings.Service {
	return s
}

func (s *FakeService) Create(ctx context.Context, query usermediabindings.CreateCommand) (usermediabindings.Model, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) DeleteOne(context.Context, usermediabindings.DeleteOneCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteMany(context.Context, usermediabindings.DeleteManyCommand) error {
	return s.ExpectedError
}

func (s *FakeService) GetByMediaID(ctx context.Context, query usermediabindings.GetByMediaIDQuery) (usermediabindings.Model, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) CheckFileBelongsToUser(ctx context.Context, query usermediabindings.CheckFileBelongsToUserQuery) (bool, error) {
	return s.ExpectedFileBelongsToUser, s.ExpectedError
}

func (s *FakeService) CheckMediaBelongsToUser(context.Context, usermediabindings.CheckMediaBelongsToUserQuery) (bool, error) {
	return s.ExpectedMediaBelongsToUser, s.ExpectedError
}

func (s *FakeService) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) ([]mediaStore.Media, error) {
	return []mediaStore.Media{}, s.ExpectedError
}
