package mediametadata

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
)

type Service interface {
	Create(context.Context, CreateCommand) (storemodels.MediaMetadataModel, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	UpdateThumbnail(context.Context, UpdateThumbnailCommand) error
}
