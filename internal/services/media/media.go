package media

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service interface {
	Create(context.Context, CreateMediaCommand) (storemodels.MediaModel, error)
	CascadeDeleteOne(context.Context, DeleteOneCommand) error
	CascadeDeleteMany(context.Context, DeleteManyCommand) (
		deletedUserMediaBindings []storemodels.UserMediaBindingsModel,
		deletedAlbumMediaBindings []storemodels.AlbumMediaBindingsModel,
		deletedMedia []storemodels.MediaModel,
		deletedMediaMetadata []storemodels.MediaMetadataModel,
		err error,
	)
	GetByUploadRequestID(context.Context, GetByUploadRequestQuery) (storemodels.MediaModel, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, GetByUploadRequestQuery) (storemodels.MediaModel, error)
	GetByFileName(context.Context, GetByFileNameQuery) (storemodels.MediaModel, error)
	GetByUserID(context.Context, GetByUserIDQuery) ([]storemodels.MediaModel, error)
	GetByMediaID(context.Context, GetByMediaIDQuery) (storemodels.MediaModel, error)
	GetByMediaIDs(context.Context, GetByMediaIDsQuery) ([]storemodels.MediaModel, error)
	GetByMediaIDsWithSort(context.Context, GetByMediaIDsWithSortQuery) ([]storemodels.MediaModel, error)
	GetTypeByFileName(context.Context, GetTypeByFileNameQuery) (string, error)
}
