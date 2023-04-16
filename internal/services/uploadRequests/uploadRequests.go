package uploadrequests

import "context"

type Service interface {
	Create(context.Context, CreateUploadRequestCommand) (Model, error)
	GetByID(context.Context, GetByIDQuery) (Model, error)
	UpdateStatus(context.Context, UpdateStatusCommand) error
}
