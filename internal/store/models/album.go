package storemodels

import (
	"github.com/rishabhkailey/media-vault/internal/constants"
	"gorm.io/gorm"
)

type AlbumModel struct {
	gorm.Model
	Name         string
	ThumbnailUrl string
	MediaCount   int
}

func (AlbumModel) TableName() string {
	return constants.ALBUMS_TABLE
}

type UserAlbumBindingsModel struct {
	gorm.Model
	UserID  string     `gorm:"index:,unique,composite:user_id_album_id"`
	AlbumID uint       `gorm:"index:,unique,composite:user_id_album_id"`
	Album   AlbumModel `gorm:"foreignKey:AlbumID"`
}

func (UserAlbumBindingsModel) TableName() string {
	return constants.USER_ALBUM_BINDINGS_TABLE
}

type AlbumMediaBindingsModel struct {
	gorm.Model
	AlbumID uint       `gorm:"index:,unique,composite:album_id_media_id"`
	MediaID uint       `gorm:"index:,unique,composite:album_id_media_id"`
	Media   MediaModel `gorm:"foreignKey:MediaID"`
	Album   AlbumModel `gorm:"foreignKey:AlbumID"`
}

func (AlbumMediaBindingsModel) TableName() string {
	return constants.ALBUM_MEDIA_BINDING_TABLE
}
