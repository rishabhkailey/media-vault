package media

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service interface {
	Create(context.Context, CreateMediaCommand) (storemodels.MediaModel, error)
	CascadeDeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByUploadRequestID(context.Context, GetByUploadRequestQuery) (storemodels.MediaModel, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, GetByUploadRequestQuery) (storemodels.MediaModel, error)
	GetByFileName(context.Context, GetByFileNameQuery) (storemodels.MediaModel, error)
	GetByUserID(context.Context, GetByUserIDQuery) ([]storemodels.MediaModel, error)
	GetByMediaID(context.Context, GetByMediaIDQuery) (storemodels.MediaModel, error)
	GetByMediaIDs(context.Context, GetByMediaIDsQuery) ([]storemodels.MediaModel, error)
	GetTypeByFileName(context.Context, GetTypeByFileNameQuery) (string, error)
}
