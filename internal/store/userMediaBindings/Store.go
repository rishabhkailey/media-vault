package usermediabindings

import (
	"context"

	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *Model) (uint, error)
	DeleteOne(ctx context.Context, userID string, mediaID uint) error
	// DeleteMany(ctx context.Context, userID string, mediaIDs []uint) error
	GetByMediaID(context.Context, uint) (Model, error)
	CheckFileBelongsToUser(ctx context.Context, userID, fileName string) (bool, error)
	// GetUserMedia(context.Context, usermediabindings.GetUserMediaQuery) ([]mediaStore.Media, error)
}
