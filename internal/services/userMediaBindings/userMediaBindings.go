package usermediabindings

import (
	"context"

	userMediaBindingsStore "github.com/rishabhkailey/media-service/internal/store/userMediaBindings"
)

// todo move get user media here
type Service interface {
	Create(context.Context, CreateCommand) (userMediaBindingsStore.Model, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	CheckFileBelongsToUser(context.Context, CheckFileBelongsToUserQuery) (bool, error)
	CheckMediaBelongsToUser(context.Context, CheckMediaBelongsToUserQuery) (bool, error)
}
