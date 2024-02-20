package uploadrequests

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
)

type Service interface {
	Create(context.Context, CreateUploadRequestCommand) (storemodels.UploadRequestsModel, error)
	DeleteOne(context.Context, DeleteOneCommand) error
	DeleteMany(context.Context, DeleteManyCommand) error
	GetByID(context.Context, GetByIDQuery) (storemodels.UploadRequestsModel, error)
	UpdateStatus(context.Context, UpdateStatusCommand) error
}
