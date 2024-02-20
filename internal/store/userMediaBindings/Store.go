package usermediabindings

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"

	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *storemodels.UserMediaBindingsModel) (uint, error)
	DeleteOne(ctx context.Context, userID string, mediaID uint) error
	// DeleteMany(ctx context.Context, userID string, mediaIDs []uint) error
	GetByMediaID(context.Context, uint) (storemodels.UserMediaBindingsModel, error)
	CheckFileBelongsToUser(ctx context.Context, userID, fileName string) (bool, error)
	CheckMultipleMediaBelongsToUser(ctx context.Context, userID string, mediaIDs []uint) (bool, error)
	// GetUserMedia(context.Context, usermediabindings.GetUserMediaQuery) ([]mediaStore.Media, error)
}
