package album

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"

	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(tx *gorm.DB) Store
	InsertAlbum(ctx context.Context,
		albumName string,
		thumbnailUrl string,
	) (album storemodels.AlbumModel, err error)
	UpdateAlbum(ctx context.Context,
		albumID uint,
		name *string,
		thumbnailUrl *string,
	) (album storemodels.AlbumModel, err error)
	InsertUserAlbumBindings(ctx context.Context,
		userID string,
		albumID uint,
	) (id uint, err error)
	GetByUserId(ctx context.Context,
		userID string,
		orderBy AlbumOrderBy,
		sort Sort,
		limit int,
		offset int,
	) (albums []storemodels.AlbumModel, err error)
	GetAlbumsByUserIdOrderByCreationAt(
		ctx context.Context,
		userID string,
		orderBy AlbumOrderBy,
		sort Sort,
		lastAlbumID *uint,
		limit int,
	) (albums []storemodels.AlbumModel, err error)
	GetAlbumsByUserIdOrderByUpdatedAt(
		ctx context.Context,
		userID string,
		orderBy AlbumOrderBy,
		sort Sort,
		lastAlbumID *uint,
		limit int,
	) (albums []storemodels.AlbumModel, err error)
	GetByID(ctx context.Context,
		albumID uint,
	) (result storemodels.AlbumModel, err error)
	// GetMediaByAlbumId(ctx context.Context,
	// 	albumID uint,
	// 	orderBy string,
	// 	sort string,
	// 	limit int,
	// 	offset int,
	// ) (mediaList []media.Media, err error)
	GetMediaByAlbumIdOrderByDate(ctx context.Context,
		albumID uint,
		lastMediaID *uint,
		sort Sort,
		limit int,
	) (mediaList []storemodels.AlbumMediaBindingsModel, err error)
	GetMediaByAlbumIdOrderByUploadDate(ctx context.Context,
		albumID uint,
		lastMediaID *uint,
		sort Sort,
		limit int,
	) (mediaList []storemodels.AlbumMediaBindingsModel, err error)
	GetMediaByAlbumIdOrderByAddedDate(ctx context.Context,
		albumID uint,
		lastMediaID *uint,
		sort Sort,
		limit int,
	) (mediaList []storemodels.AlbumMediaBindingsModel, err error)
	CheckAlbumBelongsToUser(ctx context.Context,
		userID string,
		albumID uint,
	) (ok bool, err error)
	AddMediaInAlbum(ctx context.Context,
		albumID uint,
		mediaIDs []uint,
	) (addedMediaIDs []uint, err error)
	RemoveMediaFromAlbum(ctx context.Context,
		albumID uint,
		mediaIDs []uint,
	) (removedMediaIDs []uint, err error)
	DeleteAlbum(ctx context.Context,
		albumID uint,
	) error
}
