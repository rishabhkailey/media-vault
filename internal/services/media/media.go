package media

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	WithTransaction(tx *gorm.DB) Service
	Create(context.Context, CreateMediaCommand) (Model, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByUploadRequestID(context.Context, GetByUploadRequestQuery) (Model, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, GetByUploadRequestQuery) (Model, error)
	GetByFileName(context.Context, GetByFileNameQuery) (Model, error)
	GetByUserID(context.Context, GetByUserIDQuery) ([]GetMediaQueryResultItem, error)
	GetByMediaID(context.Context, GetByMediaIDQuery) (Model, error)
	GetByMediaIDs(context.Context, GetByMediaIDsQuery) ([]GetMediaQueryResultItem, error)
	GetTypeByFileName(context.Context, GetTypeByFileNameQuery) (string, error)
}
