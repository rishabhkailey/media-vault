package newmediastorage

import (
	"context"

	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
)

type Store interface {
	SaveFile(context.Context, mediastorage.StoreSaveFileCmd) (int64, error)
	GetByFileName(ctx context.Context, fileNmae string) (mediastorage.File, error)
	DeleteOne(ctx context.Context, fileName string) error
}
