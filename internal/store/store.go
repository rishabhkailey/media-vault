package store

import (
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/store/album"
	albumstoreimpl "github.com/rishabhkailey/media-service/internal/store/album/impl"
	"github.com/rishabhkailey/media-service/internal/store/media"
	mediastoreimpl "github.com/rishabhkailey/media-service/internal/store/media/impl"
	userinfo "github.com/rishabhkailey/media-service/internal/store/userInfo"
	userinfoimpl "github.com/rishabhkailey/media-service/internal/store/userInfo/impl"
	"gorm.io/gorm"
)

// each service will have all the stores available
type Store struct {
	db            *gorm.DB
	AlbumStore    album.Store
	UserInfoStore userinfo.Store
	MediaStore    media.Store
}

func NewStore(db *gorm.DB, cache *redis.Client, minio *minio.Client) (*Store, error) {
	albumStore, err := albumstoreimpl.NewSqlStore(db, cache)
	if err != nil {
		return nil, err
	}
	userInfoStore, err := userinfoimpl.NewSqlStore(db, cache)
	if err != nil {
		return nil, err
	}
	mediaStore, err := mediastoreimpl.NewSqlStoreWithMigrate(db, cache)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:            db,
		AlbumStore:    albumStore,
		UserInfoStore: userInfoStore,
		MediaStore:    mediaStore,
	}, nil
}

func (s *Store) CreateTransaction() *gorm.DB {
	return s.db.Begin()
}
