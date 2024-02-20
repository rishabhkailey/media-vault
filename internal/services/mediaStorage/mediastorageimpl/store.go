package mediastorageimpl

import (
	"context"

	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
)

type store interface {
	SaveFile(context.Context, mediastorage.StoreSaveFileCmd) (int64, error)
	GetByFileName(ctx context.Context, fileName string) (mediastorage.File, error)
	DeleteOne(ctx context.Context, fileName string) error
}
