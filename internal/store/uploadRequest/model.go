package uploadrequest

import (
	"github.com/rishabhkailey/media-service/internal/constants"
	"gorm.io/gorm"
)

type Status string

const (
	COMPLETED_UPLOAD_STATUS   Status = "completed"
	FAILED_UPLOAD_STATUS      Status = "failed"
	IN_PROGRESS_UPLOAD_STATUS Status = "inProgress"
)

type UploadRequest struct {
	gorm.Model
	UserID string `gorm:"index:,composite:user_id_status"`
	Status Status `gorm:"index:,composite:user_id_status"`
}

func (UploadRequest) TableName() string {
	return constants.UPLOAD_REQUESTS_TABLE
}
