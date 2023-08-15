package uploadrequest

import (
	"gorm.io/gorm"
)

type Status string

const (
	COMPLETED_UPLOAD_STATUS   Status = "completed"
	FAILED_UPLOAD_STATUS      Status = "failed"
	IN_PROGRESS_UPLOAD_STATUS Status = "inProgress"
	TABLE_NAME                       = "upload_requests"
)

type UploadRequest struct {
	gorm.Model
	UserID string `gorm:"index:,composite:user_id_status"`
	Status Status `gorm:"index:,composite:user_id_status"`
}

func (UploadRequest) TableName() string {
	return TABLE_NAME
}
