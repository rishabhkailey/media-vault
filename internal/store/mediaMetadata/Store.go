package mediametadata

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *storemodels.MediaMetadataModel) (uint, error)
	DeleteOne(context.Context, uint) error
	DeleteMany(context.Context, []uint) error
	UpdateThumbnail(ctx context.Context, id uint, hasThumbnail bool, thumbnailAspectRatio float32) error
}
