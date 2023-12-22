package album

import (
	"context"

	albumStore "github.com/rishabhkailey/media-service/internal/store/album"
)

type Service interface {
	Create(context.Context, CreateAlbumCmd) (albumStore.Album, error)
	GetUserAlbums(context.Context, GetUserAlbumsQuery) ([]albumStore.Album, error)
	GetUserAlbum(context.Context, GetUserAlbumQuery) (albumStore.Album, error)
	UpdateAlbum(context.Context, UpdateAlbumCmd) (albumStore.Album, error)
	AddMedia(context.Context, AddMediaQuery) (addedMediaIDs []uint, err error)
	RemoveAlbumMedia(context.Context, RemoveMediaCmd) (removedMediaIDs []uint, err error)
	DeleteAlbum(context.Context, DeleteAlbumCmd) error
	// album will not be loaded only media will be present
	GetAlbumMedia(context.Context, GetAlbumMediaQuery) ([]albumStore.AlbumMediaBindings, error)
}
