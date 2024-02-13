package usermediabindingsimpl

import (
	"context"

	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	"gorm.io/gorm"
)

type FakeService struct {
	ExpectedError              error
	ExpectedUserMediaBinding   storemodels.UserMediaBindingsModel
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

func (s *FakeService) Create(ctx context.Context, query usermediabindings.CreateCommand) (storemodels.UserMediaBindingsModel, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) DeleteOne(context.Context, usermediabindings.DeleteOneCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteMany(context.Context, usermediabindings.DeleteManyCommand) error {
	return s.ExpectedError
}

func (s *FakeService) GetByMediaID(ctx context.Context, query usermediabindings.GetByMediaIDQuery) (storemodels.UserMediaBindingsModel, error) {
	return s.ExpectedUserMediaBinding, s.ExpectedError
}

func (s *FakeService) CheckFileBelongsToUser(ctx context.Context, query usermediabindings.CheckFileBelongsToUserQuery) (bool, error) {
	return s.ExpectedFileBelongsToUser, s.ExpectedError
}

func (s *FakeService) CheckMediaBelongsToUser(context.Context, usermediabindings.CheckMediaBelongsToUserQuery) (bool, error) {
	return s.ExpectedMediaBelongsToUser, s.ExpectedError
}

func (s *FakeService) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) ([]storemodels.MediaModel, error) {
	return []storemodels.MediaModel{}, s.ExpectedError
}
func (s *FakeService) CheckMultipleMediaBelongsToUser(ctx context.Context, query usermediabindings.CheckMultipleMediaBelongsToUserQuery) (bool, error) {
	return s.ExpectedMediaBelongsToUser, s.ExpectedError
}
