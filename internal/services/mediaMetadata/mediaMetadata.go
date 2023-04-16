package mediametadata

import "context"

type Service interface {
	Create(context.Context, CreateCommand) (Model, error)
	UpdateThumbnail(context.Context, UpdateThumbnailCommand) error
}
