package store

import (
	"github.com/go-redis/redis/v8"
	"github.com/rishabhkailey/media-service/internal/store/album"
	albumstoreimpl "github.com/rishabhkailey/media-service/internal/store/album/impl"
	"gorm.io/gorm"
)

// each service will have all the stores available
type Store struct {
	db         *gorm.DB
	AlbumStore album.Store
}

func NewStore(db *gorm.DB, cache *redis.Client) (*Store, error) {
	albumStore, err := albumstoreimpl.NewSqlStore(db, cache)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:         db,
		AlbumStore: albumStore,
	}, nil
}

func (s *Store) CreateTransaction() *gorm.DB {
	return s.db.Begin()
}
