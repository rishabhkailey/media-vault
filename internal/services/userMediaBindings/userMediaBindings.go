package usermediabindings

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
)

// todo move get user media here
type Service interface {
	Create(context.Context, CreateCommand) (storemodels.UserMediaBindingsModel, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	CheckFileBelongsToUser(context.Context, CheckFileBelongsToUserQuery) (bool, error)
	CheckMediaBelongsToUser(context.Context, CheckMediaBelongsToUserQuery) (bool, error)
	CheckMultipleMediaBelongsToUser(context.Context, CheckMultipleMediaBelongsToUserQuery) (bool, error)
}
