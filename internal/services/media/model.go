package media

import (
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

type (
	Order string
	Sort  string
)

const (
	ORDER_BY_UPLOAD_TIME         Order = "created_at"
	ORDER_BY_MEDIA_CREATION_TIME Order = "date"
	SORT_ASCENDING               Sort  = "asc"
	SORT_DESCENDING              Sort  = "desc"
	TABLE_NAME                         = "media"
	MAX_PER_PAGE_VALUE           int64 = 100
)

var (
	SUPPORTED_ORDER_BY = []Order{ORDER_BY_UPLOAD_TIME, ORDER_BY_MEDIA_CREATION_TIME}
	SUPPORTED_SORT_BY  = []Sort{SORT_ASCENDING, SORT_DESCENDING}
)

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

type GetTypeByFileNameQuery struct {
	FileName string
}

type GetByUserIDQuery struct {
	UserID  string
	OrderBy Order `form:"order" json:"order,omitempty" binding:"required"`
	Sort    Sort  `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64 `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64 `form:"perPage" json:"perPage,omitempty" binding:"required"`
}

type GetMediaQueryResultItem struct {
	MediaUrl     string `json:"url"`
	ThumbnailUrl string `json:"thumbnail_url"`
	mediametadata.Metadata
}
