package mediametadataimpl

import (
	"context"

	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	"github.com/rishabhkailey/media-service/internal/store"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service struct {
	store store.Store
}

var _ mediametadata.Service = (*Service)(nil)

func NewService(store store.Store) (mediametadata.Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd mediametadata.CreateCommand) (storemodels.MediaMetadataModel, error) {
	mediaMetadata := storemodels.MediaMetadataModel{
		MediaMetadata: cmd.MediaMetadata,
	}

	_, err := s.store.MediaMetadata.Insert(ctx, &mediaMetadata)
	return mediaMetadata, err
}

func (s *Service) DeleteOne(ctx context.Context, cmd mediametadata.DeleteOneCommand) error {
	return s.store.MediaMetadata.DeleteOne(ctx, cmd.ID)
}

func (s *Service) DeleteMany(ctx context.Context, cmd mediametadata.DeleteManyCommand) error {
	return s.store.MediaMetadata.DeleteMany(ctx, cmd.IDs)
}

func (s *Service) UpdateThumbnail(ctx context.Context, cmd mediametadata.UpdateThumbnailCommand) error {
	return s.store.MediaMetadata.UpdateThumbnail(ctx, cmd.ID, cmd.Thumbnail, cmd.ThumbnailAspectRatio)
}
