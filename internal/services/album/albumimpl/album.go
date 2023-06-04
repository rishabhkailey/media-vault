package albumimpl

import (
	"fmt"

	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/rishabhkailey/media-service/internal/services/album"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	"github.com/rishabhkailey/media-service/internal/store"
	albumStore "github.com/rishabhkailey/media-service/internal/store/album"
	"golang.org/x/net/context"
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

func (s *Service) Create(ctx context.Context, cmd album.CreateAlbumCmd) (album albumStore.Album, err error) {
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

func (s *Service) GetUserAlbums(ctx context.Context, query album.GetUserAlbumsQuery) ([]albumStore.Album, error) {
	return s.store.AlbumStore.GetByUserId(
		ctx, query.UserID,
		album.AlbumOrderAttributesMapping[query.OrderBy],
		album.AlbumSortKeywordMapping[query.Sort],
		int(query.PerPage),
		int((query.Page-1)*query.PerPage),
	)
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

func (s *Service) GetAlbumMedia(ctx context.Context, query album.GetAlbumMediaQuery) (mediaListResponse []media.GetMediaQueryResultItem, err error) {
	ok, err := s.store.AlbumStore.CheckAlbumBelongsToUser(ctx, query.UserID, query.AlbumID)
	if err != nil {
		return
	}
	if !ok {
		err = internalErrors.ErrForbidden
		return
	}
	mediaList, err := s.store.AlbumStore.GetMediaByAlbumId(
		ctx,
		query.AlbumID,
		album.AlbumMediaOrderAttributesMapping[query.OrderBy],
		album.AlbumMediaSortKeywordMapping[query.Sort],
		int(query.PerPage),
		int((query.Page-1)*query.PerPage),
	)
	if err != nil {
		return
	}
	mediaListResponse, err = mediaimpl.NewGetMediaQueryResult(mediaList)
	return
}
