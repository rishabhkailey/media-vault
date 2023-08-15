package mediaimpl

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/store"
	mediaStore "github.com/rishabhkailey/media-service/internal/store/media"
)

type Service struct {
	store store.Store
}

var _ media.Service = (*Service)(nil)

func NewService(store store.Store) (media.Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd media.CreateMediaCommand) (mediaStore.Media, error) {
	media := mediaStore.Media{
		FileName:        uuid.New().String(),
		UploadRequestID: cmd.UploadRequestID,
		MetadataID:      cmd.MetadataID,
	}
	_, err := s.store.MediaStore.Insert(ctx, &media)
	if err != nil {
		return media, fmt.Errorf("[mediaService.Create] failed: %w", err)
	}
	return media, nil
}

func (s *Service) DeleteOne(ctx context.Context, query media.DeleteOneCommand) error {
	return s.store.MediaStore.DeleteOne(ctx, query.ID)
}

func (s *Service) DeleteMany(ctx context.Context, query media.DeleteManyCommand) error {
	return s.store.MediaStore.DeleteMany(ctx, query.IDs)
}

func (s *Service) GetByUploadRequestID(ctx context.Context, query media.GetByUploadRequestQuery) (mediaStore.Media, error) {
	return s.store.MediaStore.GetByUploadRequestID(ctx, query.UploadRequestID)
}

func (s *Service) GetMediaWithMetadataByUploadRequestID(ctx context.Context, query media.GetByUploadRequestQuery) (mediaStore.Media, error) {
	return s.store.MediaStore.GetMediaWithMetadataByUploadRequestID(ctx, query.UploadRequestID)
}

func (s *Service) GetByFileName(ctx context.Context, query media.GetByFileNameQuery) (mediaStore.Media, error) {
	return s.store.MediaStore.GetByFileName(ctx, query.FileName)
}

func (s *Service) GetByUserID(ctx context.Context, query media.GetByUserIDQuery) (result []mediaStore.Media, err error) {
	if query.OrderBy == "uploaded_at" {
		return s.store.MediaStore.GetByUserIDOrderByUploadDate(ctx, query.UserID, query.LastMediaID, query.LastDate, media.SortKeywordMapping[query.Sort], int(query.PerPage))
	}
	return s.store.MediaStore.GetByUserIDOrderByDate(ctx, query.UserID, query.LastMediaID, query.LastDate, media.SortKeywordMapping[query.Sort], int(query.PerPage))
}

func (s *Service) GetTypeByFileName(ctx context.Context, query media.GetTypeByFileNameQuery) (string, error) {
	return s.store.MediaStore.GetTypeByFileName(ctx, query.FileName)
}

func (s *Service) GetByMediaIDs(ctx context.Context, query media.GetByMediaIDsQuery) (result []mediaStore.Media, err error) {
	// todo check before directly refrencing the maps
	return s.store.MediaStore.GetByMediaIDs(ctx, media.OrderKeywordMapping[query.OrderBy], media.SortKeywordMapping[query.Sort], query.MediaIDs)
}
func (s *Service) GetByMediaID(ctx context.Context, query media.GetByMediaIDQuery) (mediaStore.Media, error) {
	return s.store.MediaStore.GetByMediaID(ctx, query.MediaID)
}

func (s *Service) GetUserMediaByID(ctx context.Context, query media.UserMediaByIDQuery) (media mediaStore.Media, err error) {
	return
}
