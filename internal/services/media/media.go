package media

import "context"

type Service interface {
	Create(context.Context, CreateMediaCommand) (Model, error)
	GetByUploadRequestID(context.Context, GetByUploadRequestQuery) (Model, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, GetByUploadRequestQuery) (Model, error)
	GetByFileName(context.Context, GetByFileNameQuery) (Model, error)
	GetByUserID(context.Context, GetByUserIDQuery) ([]GetMediaQueryResultItem, error)
	GetByMediaIDs(context.Context, []uint) ([]GetMediaQueryResultItem, error)
	GetTypeByFileName(context.Context, GetTypeByFileNameQuery) (string, error)
}
