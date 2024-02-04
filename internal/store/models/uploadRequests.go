package models

import (
	"time"

	"github.com/rishabhkailey/media-service/internal/constants"
	"gorm.io/gorm"
)

type UploadRequestsModel struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index:,composite:user_id_status"`
	Status    string `gorm:"index:,composite:user_id_status"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UploadRequestsModel) TableName() string {
	return constants.UPLOAD_REQUESTS_TABLE
}
