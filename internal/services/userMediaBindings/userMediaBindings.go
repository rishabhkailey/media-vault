package usermediabindings

import (
	"context"

	mediaStore "github.com/rishabhkailey/media-service/internal/store/media"
	"gorm.io/gorm"
)

// todo move get user media here
type Service interface {
	WithTransaction(tx *gorm.DB) Service
	Create(context.Context, CreateCommand) (Model, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByMediaID(context.Context, GetByMediaIDQuery) (Model, error)
	CheckFileBelongsToUser(context.Context, CheckFileBelongsToUserQuery) (bool, error)
	CheckMediaBelongsToUser(context.Context, CheckMediaBelongsToUserQuery) (bool, error)
	GetUserMedia(context.Context, GetUserMediaQuery) ([]mediaStore.Media, error)
}
