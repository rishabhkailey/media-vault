package usermediabindings

import (
	"github.com/rishabhkailey/media-service/internal/services/media"
	"gorm.io/gorm"
)

const (
	ORDER_BY_UPLOAD_TIME         = "created_at"
	ORDER_BY_MEDIA_CREATION_TIME = "date"
	SORT_ASCENDING               = "asc"
	SORT_DESCENDING              = "desc"
	TABLE_NAME                   = "user_media_bindings"
)

var SUPPORTED_ORDER_BY = []string{ORDER_BY_UPLOAD_TIME, ORDER_BY_MEDIA_CREATION_TIME}

// todo user subject as userid on the client so even if user change email or user name the user should not loose its data
type Model struct {
	gorm.Model
	UserID  string      `gorm:"index:,unique,composite:user_id_media_id"`
	MediaID uint        `gorm:"index:,unique,composite:user_id_media_id"`
	Media   media.Model `gorm:"foreignKey:MediaID"`
}

func (Model) TableName() string {
	return TABLE_NAME
}

type CreateCommand struct {
	UserID  string
	MediaID uint
}

type GetByMediaIDQuery struct {
	MediaID uint
}

type CheckFileBelongsToUserQuery struct {
	UserID   string
	FileName string
}

type GetUserMediaQuery struct {
	UserID  string
	OrderBy string
	Sort    string
	Offset  int
	Limit   int
}
