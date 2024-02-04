package uploadrequests

import (
	"time"

	"github.com/rishabhkailey/media-service/internal/constants"
	"gorm.io/gorm"
)

type Status string

const (
	COMPLETED_UPLOAD_STATUS   Status = "completed"
	FAILED_UPLOAD_STATUS      Status = "failed"
	IN_PROGRESS_UPLOAD_STATUS Status = "inProgress"
)

type Model struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index:,composite:user_id_status"`
	Status    Status `gorm:"index:,composite:user_id_status"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Model) TableName() string {
	return constants.UPLOAD_REQUESTS_TABLE
}

// will be used by upload request service
type CreateUploadRequestCommand struct {
	UserID string
}

type GetByIDQuery struct {
	ID string
}

type UpdateStatusCommand struct {
	Status Status
	ID     string
}

type DeleteOneCommand struct {
	ID string
}

type DeleteManyCommand struct {
	IDs []string
}
