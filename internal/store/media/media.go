package media

import (
	"context"
	"time"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *storemodels.MediaModel) (uint, error)
	// deletes usermediabinding, medialabumbinding, media and media metadata
	CascadeDeleteOne(ctx context.Context, mediaID uint, userID string, mediaMetadataId uint) error
	CascadeDeleteMany(ctx context.Context, userID string, mediaIDs []uint) (
		deletedUserMediaBindings []storemodels.UserMediaBindingsModel,
		deletedAlbumMediaBindings []storemodels.AlbumMediaBindingsModel,
		deletedMedia []storemodels.MediaModel,
		deletedMediaMetadata []storemodels.MediaMetadataModel,
		err error,
	)
	DeleteMany(context.Context, []uint) error
	GetByUploadRequestID(context.Context, string) (storemodels.MediaModel, error)
	GetByFileName(context.Context, string) (storemodels.MediaModel, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, string) (storemodels.MediaModel, error)
	GetByUserIDOrderByDate(ctx context.Context, userID string, lastMediaID *uint, lastDate *time.Time, sort Sort, limit int) ([]storemodels.MediaModel, error)
	GetByUserIDOrderByUploadDate(ctx context.Context, userID string, lastMediaID *uint, lastDate *time.Time, sort Sort, limit int) ([]storemodels.MediaModel, error)
	GetByMediaID(ctx context.Context, mediaID uint) (storemodels.MediaModel, error)
	GetByMediaIDs(ctx context.Context, mediaIDs []uint) ([]storemodels.MediaModel, error)
	GetByMediaIDsWithSort(ctx context.Context, orderBy OrderBy, sort Sort, mediaIDs []uint) ([]storemodels.MediaModel, error)
	GetTypeByFileName(context.Context, string) (string, error)
}
