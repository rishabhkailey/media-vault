package album

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service interface {
	Create(context.Context, CreateAlbumCmd) (storemodels.AlbumModel, error)
	GetUserAlbums(context.Context, GetUserAlbumsQuery) ([]storemodels.AlbumModel, error)
	GetUserAlbum(context.Context, GetUserAlbumQuery) (storemodels.AlbumModel, error)
	UpdateAlbum(context.Context, UpdateAlbumCmd) (storemodels.AlbumModel, error)
	AddMedia(context.Context, AddMediaQuery) (addedMediaIDs []uint, err error)
	RemoveAlbumMedia(context.Context, RemoveMediaCmd) (removedMediaIDs []uint, err error)
	DeleteAlbum(context.Context, DeleteAlbumCmd) error
	// album will not be loaded only media will be present
	GetAlbumMedia(context.Context, GetAlbumMediaQuery) ([]storemodels.AlbumMediaBindingsModel, error)
}
