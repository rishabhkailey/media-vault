package uploadrequests

import (
	"context"

	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"gorm.io/gorm"
)

type Store interface {
	WithTransaction(*gorm.DB) Store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *storemodels.UploadRequestsModel) (string, error)
	DeleteOne(context.Context, string) error
	DeleteMany(context.Context, []string) error
	GetByID(context.Context, string) (storemodels.UploadRequestsModel, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
