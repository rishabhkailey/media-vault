package album

import (
	"context"

	albumStore "github.com/rishabhkailey/media-service/internal/store/album"
	"github.com/rishabhkailey/media-service/internal/store/media"
)

type Service interface {
	Create(context.Context, CreateAlbumCmd) (albumStore.Album, error)
	GetUserAlbums(context.Context, GetUserAlbumsQuery) ([]albumStore.Album, error)
	GetUserAlbum(context.Context, GetUserAlbumQuery) (albumStore.Album, error)
	UpdateAlbum(context.Context, UpdateAlbumCmd) (albumStore.Album, error)
	AddMedia(context.Context, AddMediaQuery) (addedMediaIDs []uint, err error)
	RemoveAlbumMedia(context.Context, RemoveMediaCmd) (removedMediaIDs []uint, err error)
	DeleteAlbum(context.Context, DeleteAlbumCmd) error
	GetAlbumMedia(context.Context, GetAlbumMediaQuery) ([]media.Media, error)
}
