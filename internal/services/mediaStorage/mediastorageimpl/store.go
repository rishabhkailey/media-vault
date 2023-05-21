package mediastorageimpl

import (
	"context"

	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
)

type store interface {
	SaveFile(context.Context, mediastorage.StoreSaveFileCmd) (int64, error)
	GetByFileName(ctx context.Context, fileNmae string) (mediastorage.File, error)
	DeleteOne(ctx context.Context, fileName string) error
	// DeleteMany(ctx context.Context, fileNames []string) (failed []string, errs []error)
}
