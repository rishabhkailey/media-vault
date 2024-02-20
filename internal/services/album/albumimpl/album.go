package albumimpl

import (
	"errors"
	"fmt"

	internalErrors "github.com/rishabhkailey/media-vault/internal/errors"
	"github.com/rishabhkailey/media-vault/internal/services/album"
	"github.com/rishabhkailey/media-vault/internal/store"
	albumStore "github.com/rishabhkailey/media-vault/internal/store/album"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Service struct {
	store store.Store
}

var _ album.Service = (*Service)(nil)

func NewService(store store.Store) (album.Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd album.CreateAlbumCmd) (album storemodels.AlbumModel, err error) {
	tx := s.store.CreateTransaction()
	tx.Begin()
	album, err = s.store.AlbumStore.WithTransaction(tx).InsertAlbum(ctx, cmd.Name, cmd.ThumbnailUrl)
	if err != nil {
		err = fmt.Errorf("[Services.Album.Create] insert album failed :%w", err)
		tx.Rollback()
		return
	}
	_, err = s.store.AlbumStore.WithTransaction(tx).InsertUserAlbumBindings(ctx, cmd.UserID, album.ID)
	if err != nil {
		err = fmt.Errorf("[Services.Album.Create] insert user album binding failed :%w", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (s *Service) GetUserAlbums(ctx context.Context, query album.GetUserAlbumsQuery) ([]storemodels.AlbumModel, error) {
	switch query.OrderBy {
	case albumStore.AlbumOrderByDate:
		{
			return s.store.AlbumStore.GetAlbumsByUserIdOrderByCreationAt(
				ctx,
				query.UserID,
				query.OrderBy,
				query.Sort,
				query.LastAlbumID,
				query.PerPage,
			)
		}
	case albumStore.AlbumOrderByUpdatedAt:
		{
			return s.store.AlbumStore.GetAlbumsByUserIdOrderByUpdatedAt(
				ctx,
				query.UserID,
				query.OrderBy,
				query.Sort,
				query.LastAlbumID,
				query.PerPage,
			)
		}
	default:
		{
			return s.store.AlbumStore.GetAlbumsByUserIdOrderByCreationAt(
				ctx,
				query.UserID,
				query.OrderBy,
				query.Sort,
				query.LastAlbumID,
				query.PerPage,
			)
		}
	}
}

func (s *Service) GetUserAlbum(ctx context.Context, query album.GetUserAlbumQuery) (album storemodels.AlbumModel, err error) {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, query.UserID, query.AlbumID)
	if err != nil {
		return
	}
	if !ok {
		err = internalErrors.ErrForbidden
		return
	}
	album, err = s.store.AlbumStore.GetByID(ctx, query.AlbumID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = internalErrors.ErrForbidden
	}
	return
}

func (s *Service) AddMedia(ctx context.Context, query album.AddMediaQuery) (addedMediaIDs []uint, err error) {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, query.UserID, query.AlbumID)
	if err != nil {
		return
	}
	if !ok {
		err = internalErrors.ErrForbidden
		return
	}
	return s.store.AlbumStore.AddMediaInAlbum(ctx, query.AlbumID, query.MediaIDs)
}

func (s *Service) RemoveAlbumMedia(ctx context.Context, query album.RemoveMediaCmd) ([]uint, error) {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, query.UserID, query.AlbumID)
	if err != nil {
		return []uint{}, err
	}
	if !ok {
		return []uint{}, internalErrors.ErrForbidden
	}
	return s.store.AlbumStore.RemoveMediaFromAlbum(ctx, query.AlbumID, query.MediaIDs)
}

func (s *Service) DeleteAlbum(ctx context.Context, cmd album.DeleteAlbumCmd) error {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, cmd.UserID, cmd.AlbumID)
	if err != nil {
		return err
	}
	if !ok {
		return internalErrors.ErrForbidden
	}
	return s.store.AlbumStore.DeleteAlbum(ctx, cmd.AlbumID)
}

func (s *Service) GetAlbumMedia(ctx context.Context, query album.GetAlbumMediaQuery) (albumMediaBindings []storemodels.AlbumMediaBindingsModel, err error) {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, query.UserID, query.AlbumID)
	if err != nil {
		return
	}
	if !ok {
		err = internalErrors.ErrForbidden
		return
	}
	switch query.OrderBy {
	case albumStore.AlbumMediaOrderByAddedDate:
		return s.store.AlbumStore.GetMediaByAlbumIdOrderByAddedDate(
			ctx,
			query.AlbumID,
			query.LastMediaID,
			query.Sort,
			int(query.PerPage),
		)
	case albumStore.AlbumMediaOrderByMediaDate:
		return s.store.AlbumStore.GetMediaByAlbumIdOrderByDate(
			ctx,
			query.AlbumID,
			query.LastMediaID,
			query.Sort,
			int(query.PerPage),
		)
	case albumStore.AlbumMediaOrderByUploadedDate:
		return s.store.AlbumStore.GetMediaByAlbumIdOrderByUploadDate(
			ctx,
			query.AlbumID,
			query.LastMediaID,
			query.Sort,
			int(query.PerPage),
		)
	default:
		return s.store.AlbumStore.GetMediaByAlbumIdOrderByUploadDate(
			ctx,
			query.AlbumID,
			query.LastMediaID,
			query.Sort,
			int(query.PerPage),
		)
	}

}

func (s *Service) UpdateAlbum(ctx context.Context, cmd album.UpdateAlbumCmd) (storemodels.AlbumModel, error) {
	return s.store.AlbumStore.UpdateAlbum(ctx, cmd.AlbumID, cmd.Name, cmd.ThumbnailUrl)
}
