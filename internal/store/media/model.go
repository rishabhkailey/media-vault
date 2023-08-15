package media

import (
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

type OrderBy string
type Sort string

const (
	TABLE_NAME         = "media"
	Date       OrderBy = "date"
	UploadedAt OrderBy = "created_at"
	Ascending  Sort    = "asc"
	Descending Sort    = "desc"
)

type Media struct {
	gorm.Model
	FileName        string `gorm:"index,unique"`
	UploadRequestID string `gorm:"index"`
	UploadRequest   uploadrequests.Model
	Metadata        mediametadata.Model
	MetadataID      uint
}

func (Media) TableName() string {
	return TABLE_NAME
}
