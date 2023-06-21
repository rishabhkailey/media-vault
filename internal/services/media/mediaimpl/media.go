package mediaimpl

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"gorm.io/gorm"
)

type Service struct {
	store store
}

var _ media.Service = (*Service)(nil)

func NewService(db *gorm.DB, cache *redis.Client) (media.Service, error) {
	store, err := newSqlStoreWithMigrate(db, cache)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) WithTransaction(tx *gorm.DB) media.Service {
	return &Service{
		store: s.store.WithTransaction(tx),
	}
}

func (s *Service) Create(ctx context.Context, cmd media.CreateMediaCommand) (media.Model, error) {
	media := media.Model{
		FileName:        uuid.New().String(),
		UploadRequestID: cmd.UploadRequestID,
		MetadataID:      cmd.MetadataID,
	}
	_, err := s.store.Insert(ctx, &media)
	if err != nil {
		return media, fmt.Errorf("[mediaService.Create] failed: %w", err)
	}
	return media, nil
}

func (s *Service) DeleteOne(ctx context.Context, query media.DeleteOneCommand) error {
	return s.store.DeleteOne(ctx, query.ID)
}

func (s *Service) DeleteMany(ctx context.Context, query media.DeleteManyCommand) error {
	return s.store.DeleteMany(ctx, query.IDs)
}

func (s *Service) GetByUploadRequestID(ctx context.Context, query media.GetByUploadRequestQuery) (media.Model, error) {
	return s.store.GetByUploadRequestID(ctx, query.UploadRequestID)
}

func (s *Service) GetMediaWithMetadataByUploadRequestID(ctx context.Context, query media.GetByUploadRequestQuery) (media.Model, error) {
	return s.store.GetMediaWithMetadataByUploadRequestID(ctx, query.UploadRequestID)
}

func (s *Service) GetByFileName(ctx context.Context, query media.GetByFileNameQuery) (media.Model, error) {
	return s.store.GetByFileName(ctx, query.FileName)
}

func (s *Service) GetByUserID(ctx context.Context, query media.GetByUserIDQuery) (result []media.Model, err error) {
	if query.OrderBy == "uploaded_at" {
		return s.store.GetByUserIDOrderByUploadDate(ctx, query.UserID, query.LastMediaID, query.LastDate, media.SortKeywordMapping[query.Sort], int(query.PerPage))
	}
	return s.store.GetByUserIDOrderByDate(ctx, query.UserID, query.LastMediaID, query.LastDate, media.SortKeywordMapping[query.Sort], int(query.PerPage))
}

func (s *Service) GetTypeByFileName(ctx context.Context, query media.GetTypeByFileNameQuery) (string, error) {
	return s.store.GetTypeByFileName(ctx, query.FileName)
}

func (s *Service) GetByMediaIDs(ctx context.Context, query media.GetByMediaIDsQuery) (result []media.Model, err error) {
	return s.store.GetByMediaIDs(ctx, query)
}
func (s *Service) GetByMediaID(ctx context.Context, query media.GetByMediaIDQuery) (media.Model, error) {
	return s.store.GetByMediaID(ctx, query)
}
