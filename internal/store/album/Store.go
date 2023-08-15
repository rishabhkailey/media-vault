package album

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/store/media"
	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(tx *gorm.DB) Store
	InsertAlbum(ctx context.Context, albumName string, thumbnailUrl string) (album Album, err error)
	UpdateAlbum(ctx context.Context, albumID uint, name *string, thumbnailUrl *string) (album Album, err error)
	InsertUserAlbumBindings(ctx context.Context, userID string, albumID uint) (id uint, err error)
	GetByUserId(ctx context.Context, userID string, orderBy string, sort string, limit int, offset int) (albums []Album, err error)
	GetByID(ctx context.Context, albumID uint) (result Album, err error)
	GetMediaByAlbumId(ctx context.Context, albumID uint, orderBy string, sort string, limit int, offset int) (mediaList []media.Media, err error)
	CheckAlbumBelongsToUser(ctx context.Context, userID string, albumID uint) (ok bool, err error)
	AddMediaInAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (addedMediaIDs []uint, err error)
	RemoveMediaFromAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (removedMediaIDs []uint, err error)
	DeleteAlbum(ctx context.Context, albumID uint) error
	UpdateThumbnail(ctx context.Context, mediaID uint, thumbnail bool, thumbnailAspectRatio float32) error
}
