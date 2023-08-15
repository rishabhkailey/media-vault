package uploadrequestimpl

import (
	"github.com/go-redis/redis/v8"
	media "github.com/rishabhkailey/media-service/internal/store/media"
	uploadrequest "github.com/rishabhkailey/media-service/internal/store/uploadRequest"
	"gorm.io/gorm"
)

// todo cache should have different store?
// maybe we can do something similar to media storage(cache wrapper for store)
type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

var _ uploadrequest.Store = (*sqlStore)(nil)

// todo
// func newSqlStore(db *gorm.DB, cache *redis.Client) store {
// 	return &sqlStore{
// 		db:    db,
// 		cache: cache,
// 	}
// }

func NewSqlStoreWithMigrate(db *gorm.DB, cache *redis.Client) (uploadrequest.Store, error) {
	if err := db.Migrator().AutoMigrate(&media.Media{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) uploadrequest.Store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

// err := s.db.WithContext(ctx).Create(&uploadRequest).Error
func (s *sqlStore) CreateUploadRequest(userID string) (uploadRequest uploadrequest.UploadRequest, err error) {
	return
}
