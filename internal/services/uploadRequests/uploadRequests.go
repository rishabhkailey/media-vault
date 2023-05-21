package uploadrequests

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	WithTransaction(tx *gorm.DB) Service
	Create(context.Context, CreateUploadRequestCommand) (Model, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByID(context.Context, GetByIDQuery) (Model, error)
	UpdateStatus(context.Context, UpdateStatusCommand) error
}
