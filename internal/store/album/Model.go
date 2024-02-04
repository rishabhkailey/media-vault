package album

import (
	"github.com/rishabhkailey/media-service/internal/constants"
	"github.com/rishabhkailey/media-service/internal/store/media"
	"gorm.io/gorm"
)

// AlbumOrderBy represent the attributes of Album table that can be used to sort albums
type AlbumOrderBy string

// AlbumOrderBy represent the attributes of AlbumMediaBindings store that can be used to sort media of a album
type AlbumMediaOrderBy string
type Sort string

const (
	USER_ALBUM_BINDINGS_TABLE_NAME                   = "user_album_bindings"
	ALBUM_MEDIA_BINDING_TABLE_NAME                   = "album_media_bindings"
	AlbumOrderByDate               AlbumOrderBy      = "album.created_at"
	AlbumOrderByUpdatedAt          AlbumOrderBy      = "album.updated_at"
	AlbumMediaOrderByAddedDate     AlbumMediaOrderBy = "album_media_bindings.created_at"
	AlbumMediaOrderByUploadedDate  AlbumMediaOrderBy = "media.created_at"
	AlbumMediaOrderByMediaDate     AlbumMediaOrderBy = "media_metadata.date"
	Ascending                      Sort              = "asc"
	Descending                     Sort              = "desc"
)

type Album struct {
	gorm.Model
	Name         string
	ThumbnailUrl string
	MediaCount   int
}

func (Album) TableName() string {
	return constants.ALBUMS_TABLE
}

type UserAlbumBindings struct {
	gorm.Model
	UserID  string `gorm:"index:,unique,composite:user_id_album_id"`
	AlbumID uint   `gorm:"index:,unique,composite:user_id_album_id"`
	Album   Album  `gorm:"foreignKey:AlbumID"`
}

func (UserAlbumBindings) TableName() string {
	return constants.USER_ALBUM_BINDINGS_TABLE
}

type AlbumMediaBindings struct {
	gorm.Model
	AlbumID uint        `gorm:"index:,unique,composite:album_id_media_id"`
	MediaID uint        `gorm:"index:,unique,composite:album_id_media_id"`
	Media   media.Media `gorm:"foreignKey:MediaID"`
	Album   Album       `gorm:"foreignKey:AlbumID"`
}

func (AlbumMediaBindings) TableName() string {
	return constants.ALBUM_MEDIA_BINDING_TABLE
}
