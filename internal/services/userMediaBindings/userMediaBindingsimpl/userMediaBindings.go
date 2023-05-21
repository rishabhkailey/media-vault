package usermediabindingsimpl

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"gorm.io/gorm"
)

type Service struct {
	store store
}

var _ usermediabindings.Service = (*Service)(nil)

func NewService(db *gorm.DB) (usermediabindings.Service, error) {
	store, err := newSqlStore(db)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) WithTransaction(tx *gorm.DB) usermediabindings.Service {
	return &Service{
		store: s.store.WithTransaction(tx),
	}
}

func (s *Service) Create(ctx context.Context, cmd usermediabindings.CreateCommand) (usermediabindings.Model, error) {
	userMediaBinding := usermediabindings.Model{
		UserID:  cmd.UserID,
		MediaID: cmd.MediaID,
	}

	_, err := s.store.Insert(ctx, &userMediaBinding)
	return userMediaBinding, err
}

func (s *Service) DeleteOne(ctx context.Context, cmd usermediabindings.DeleteOneCommand) error {
	return s.store.DeleteOne(ctx, cmd.UserID, cmd.MediaID)
}

func (s *Service) DeleteMany(ctx context.Context, cmd usermediabindings.DeleteManyCommand) error {
	return s.store.DeleteMany(ctx, cmd.UserID, cmd.MediaIDs)
}

func (s *Service) GetByMediaID(ctx context.Context, query usermediabindings.GetByMediaIDQuery) (usermediabindings.Model, error) {
	return s.store.GetByMediaID(ctx, query.MediaID)
}

func (s *Service) CheckFileBelongsToUser(ctx context.Context, cmd usermediabindings.CheckFileBelongsToUserQuery) (bool, error) {
	return s.store.CheckFileBelongsToUser(ctx, cmd)
}

func (s *Service) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) (mediaList []media.Model, err error) {
	return s.store.GetUserMedia(ctx, query)
}

func (s *Service) CheckMediaBelongsToUser(ctx context.Context, query usermediabindings.CheckMediaBelongsToUserQuery) (bool, error) {
	userMediaBinding, err := s.store.GetByMediaID(ctx, query.MediaID)
	if err != nil {
		return false, err
	}
	return userMediaBinding.UserID == query.UserID, err
}
