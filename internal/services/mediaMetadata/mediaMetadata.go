package mediametadata

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	WithTransaction(tx *gorm.DB) Service
	Create(context.Context, CreateCommand) (Model, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	UpdateThumbnail(context.Context, UpdateThumbnailCommand) error
}
