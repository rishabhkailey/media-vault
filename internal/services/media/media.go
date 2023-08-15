package media

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/store/media"
)

type Service interface {
	Create(context.Context, CreateMediaCommand) (media.Media, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByUploadRequestID(context.Context, GetByUploadRequestQuery) (media.Media, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, GetByUploadRequestQuery) (media.Media, error)
	GetByFileName(context.Context, GetByFileNameQuery) (media.Media, error)
	GetByUserID(context.Context, GetByUserIDQuery) ([]media.Media, error)
	GetByMediaID(context.Context, GetByMediaIDQuery) (media.Media, error)
	GetUserMediaByID(context.Context, UserMediaByIDQuery) (media.Media, error)
	GetByMediaIDs(context.Context, GetByMediaIDsQuery) ([]media.Media, error)
	GetTypeByFileName(context.Context, GetTypeByFileNameQuery) (string, error)
}
