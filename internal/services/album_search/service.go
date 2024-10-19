package albumsearch

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
)

type Service interface {
	GetAlbumIDs(context.Context, AlbumSearchQuery) ([]uint, error)
	CreateOne(ctx context.Context, album storemodels.AlbumModel, userID string) (int64, error)
	CreateMany(ctx context.Context, albums []storemodels.AlbumModel, userID string) (int64, error)
	DeleteOne(context.Context, int64) (int64, error)
	DeleteMany(context.Context, []int64) (int64, error)
}
