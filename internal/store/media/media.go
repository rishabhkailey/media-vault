package media

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *Media) (uint, error)
	DeleteOne(context.Context, uint) error
	DeleteMany(context.Context, []uint) error
	GetByUploadRequestID(context.Context, string) (Media, error)
	GetByFileName(context.Context, string) (Media, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, string) (Media, error)
	GetByUserIDOrderByDate(ctx context.Context, userID string, lastMediaID *uint, lastDate *time.Time, sort Sort, limit int) ([]Media, error)
	GetByUserIDOrderByUploadDate(ctx context.Context, userID string, lastMediaID *uint, lastDate *time.Time, sort Sort, limit int) ([]Media, error)
	GetByMediaID(ctx context.Context, mediaID uint) (Media, error)
	GetByMediaIDs(ctx context.Context, orderBy OrderBy, sort Sort, mediaIDs []uint) ([]Media, error)
	GetTypeByFileName(context.Context, string) (string, error)
}
