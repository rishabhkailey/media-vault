package mediametadataimpl

import (
	"context"

	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	"gorm.io/gorm"
)

type Service struct {
	store store
}

var _ mediametadata.Service = (*Service)(nil)

func NewService(db *gorm.DB) (mediametadata.Service, error) {
	store, err := newSqlStore(db)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd mediametadata.CreateCommand) (mediametadata.Model, error) {
	mediaMetadata := mediametadata.Model{
		Metadata: cmd.Metadata,
	}

	_, err := s.store.Insert(ctx, &mediaMetadata)
	return mediaMetadata, err
}

func (s *Service) UpdateThumbnail(ctx context.Context, cmd mediametadata.UpdateThumbnailCommand) error {
	return s.store.UpdateThumbnail(ctx, cmd)
}
