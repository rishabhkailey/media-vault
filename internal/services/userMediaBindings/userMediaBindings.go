package usermediabindings

import (
	"context"

	"github.com/rishabhkailey/media-service/internal/services/media"
)

// todo move get user media here
type Service interface {
	Create(context.Context, CreateCommand) (Model, error)
	GetByMediaID(context.Context, GetByMediaIDQuery) (Model, error)
	CheckFileBelongsToUser(context.Context, CheckFileBelongsToUserQuery) (bool, error)
	GetUserMedia(context.Context, GetUserMediaQuery) ([]media.Model, error)
}
