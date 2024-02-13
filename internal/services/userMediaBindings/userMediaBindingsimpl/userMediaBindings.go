package usermediabindingsimpl

import (
	"context"

	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-service/internal/store"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service struct {
	store store.Store
}

var _ usermediabindings.Service = (*Service)(nil)

func NewService(store store.Store) (usermediabindings.Service, error) {
	return &Service{
		store: store,
	}, nil
}

// func (s *Service) WithTransaction(tx *gorm.DB) usermediabindings.Service {
// 	return &Service{
// 		store: s.store.WithTransaction(tx),
// 	}
// }

func (s *Service) Create(ctx context.Context, cmd usermediabindings.CreateCommand) (storemodels.UserMediaBindingsModel, error) {
	userMediaBinding := storemodels.UserMediaBindingsModel{
		UserID:  cmd.UserID,
		MediaID: cmd.MediaID,
	}

	_, err := s.store.UserMediaBindings.Insert(ctx, &userMediaBinding)
	return userMediaBinding, err
}

func (s *Service) DeleteOne(ctx context.Context, cmd usermediabindings.DeleteOneCommand) error {
	return s.store.UserMediaBindings.DeleteOne(ctx, cmd.UserID, cmd.MediaID)
}

// func (s *Service) DeleteMany(ctx context.Context, cmd usermediabindings.DeleteManyCommand) error {
// 	return s.store.DeleteMany(ctx, cmd.UserID, cmd.MediaIDs)
// }

// func (s *Service) GetByMediaID(ctx context.Context, query usermediabindings.GetByMediaIDQuery) (usermediabindings.Model, error) {
// 	return s.store.GetByMediaID(ctx, query.MediaID)
// }

func (s *Service) CheckFileBelongsToUser(ctx context.Context, cmd usermediabindings.CheckFileBelongsToUserQuery) (bool, error) {
	return s.store.UserMediaBindings.CheckFileBelongsToUser(ctx, cmd.UserID, cmd.FileName)
}

// func (s *Service) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) (mediaList []mediaStore.Media, err error) {
// 	return s.store.GetUserMedia(ctx, query)
// }

func (s *Service) CheckMediaBelongsToUser(ctx context.Context, query usermediabindings.CheckMediaBelongsToUserQuery) (bool, error) {
	userMediaBinding, err := s.store.UserMediaBindings.GetByMediaID(ctx, query.MediaID)
	if err != nil {
		return false, err
	}
	return userMediaBinding.UserID == query.UserID, err
}

func (s *Service) CheckMultipleMediaBelongsToUser(ctx context.Context, query usermediabindings.CheckMultipleMediaBelongsToUserQuery) (bool, error) {
	return s.store.UserMediaBindings.CheckMultipleMediaBelongsToUser(ctx, query.UserID, query.MediaIDs)
}
