package media

import (
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

const (
	ORDER_BY_UPLOAD_TIME         = "created_at"
	ORDER_BY_MEDIA_CREATION_TIME = "date"
	SORT_ASCENDING               = "asc"
	SORT_DESCENDING              = "desc"
	TABLE_NAME                   = "media"
)

var SUPPORTED_ORDER_BY = []string{ORDER_BY_UPLOAD_TIME, ORDER_BY_MEDIA_CREATION_TIME}

type Model struct {
	gorm.Model
	FileName        string `gorm:"index,unique"`
	UploadRequestID string `gorm:"index"`
	UploadRequest   uploadrequests.Model
	Metadata        mediametadata.Model
	MetadataID      uint
}

func (Model) TableName() string {
	return TABLE_NAME
}

// will be used by upload request service
type CreateMediaCommand struct {
	UploadRequestID string
	MetadataID      uint
}

type GetByUploadRequestQuery struct {
	UploadRequestID string
}

type GetByFileNameQuery struct {
	FileName string
}

type GetByUserIDQuery struct {
	UserID  string
	OrderBy string
	Sort    string
	Offset  int
	Limit   int
}

type GetTypeByFileNameQuery struct {
	FileName string
}
