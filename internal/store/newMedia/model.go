package newmedia

import (
	"github.com/rishabhkailey/media-service/internal/constants"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	"gorm.io/gorm"
)

type OrderBy string
type Sort string

const (
	Date       OrderBy = "date"
	UploadedAt OrderBy = "created_at"
	Ascending  Sort    = "asc"
	Descending Sort    = "desc"
)

type Media struct {
	gorm.Model
	FileName        string `gorm:"index,unique"`
	UploadRequestID string `gorm:"index"`
	UploadRequest   storemodels.UploadRequestsModel
	Metadata        storemodels.MediaMetadataModel
	MetadataID      uint
}

func (Media) TableName() string {
	return constants.MEDIA_TABLE
}
