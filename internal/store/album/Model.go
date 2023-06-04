package album

import (
	"github.com/rishabhkailey/media-service/internal/services/media"
	"gorm.io/gorm"
)

const (
	ALBUMS_TABLE_NAME              = "albums"
	USER_ALBUM_BINDINGS_TABLE_NAME = "user_album_bindings"
	ALBUM_MEDIA_BINDING_TABLE_NAME = "album_media_bindings"
)

type Album struct {
	gorm.Model
	Name         string
	ThumbnailUrl string
	MediaCount   int
}

func (Album) TableName() string {
	return ALBUMS_TABLE_NAME
}

type UserAlbumBindings struct {
	gorm.Model
	UserID  string             `gorm:"index:,unique,composite:user_id_album_id"`
	AlbumID uint               `gorm:"index:,unique,composite:user_id_album_id"`
	Alubm   AlbumMediaBindings `gorm:"foreignKey:AlbumID"`
}

func (UserAlbumBindings) TableName() string {
	return USER_ALBUM_BINDINGS_TABLE_NAME
}

type AlbumMediaBindings struct {
	gorm.Model
	AlbumID uint        `gorm:"index:,unique,composite:album_id_media_id"`
	MediaID uint        `gorm:"index:,unique,composite:album_id_media_id"`
	Media   media.Model `gorm:"foreignKey:MediaID"`
	Album   Album       `gorm:"foreignKey:AlbumID"`
}

func (AlbumMediaBindings) TableName() string {
	return ALBUM_MEDIA_BINDING_TABLE_NAME
}
