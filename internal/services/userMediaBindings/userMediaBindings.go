package usermediabindings

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

// todo move get user media here
type Service interface {
	Create(context.Context, CreateCommand) (storemodels.UserMediaBindingsModel, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	CheckFileBelongsToUser(context.Context, CheckFileBelongsToUserQuery) (bool, error)
	CheckMediaBelongsToUser(context.Context, CheckMediaBelongsToUserQuery) (bool, error)
}
