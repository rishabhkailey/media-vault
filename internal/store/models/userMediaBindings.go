package storemodels

import (
	"github.com/rishabhkailey/media-service/internal/constants"
	"gorm.io/gorm"
)

// todo user subject as userid on the client so even if user change email or user name the user should not loose its data
type UserMediaBindingsModel struct {
	gorm.Model
	UserID  string     `gorm:"index:,unique,composite:user_id_media_id"`
	MediaID uint       `gorm:"index:,unique,composite:user_id_media_id"`
	Media   MediaModel `gorm:"foreignKey:MediaID"`
}

func (UserMediaBindingsModel) TableName() string {
	return constants.USER_MEDIA_BINDINGS_TABLE
}
