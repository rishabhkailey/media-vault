package media

import (
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

const (
	TABLE_NAME               = "media"
	MAX_PER_PAGE_VALUE int64 = 100
)

var (
	SortKeywordMapping = map[string]string{
		"asc":        "asc",
		"desc":       "desc",
		"ascending":  "asc",
		"descending": "desc",
	}

	OrderAttributesMapping = map[string]string{
		"date":        "date",
		"uploaded_at": "created_at",
	}
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
	OrderBy string `form:"order" json:"order,omitempty" binding:"required"`
	Sort    string `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64  `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64  `form:"perPage" json:"perPage,omitempty" binding:"required"`
}

type GetMediaQueryResultItem struct {
	MediaUrl     string `json:"url"`
	ThumbnailUrl string `json:"thumbnail_url"`
	mediametadata.Metadata
}

type GetByMediaIDsQuery struct {
	MediaIDs []uint
	OrderBy  string
	Sort     string
}
