package userinfoimpl

import (
	"context"

	"github.com/go-redis/redis/v8"
	userinfo "github.com/rishabhkailey/media-vault/internal/store/userInfo"
	"gorm.io/gorm"
)

type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) userinfo.Store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

var _ userinfo.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB, cache *redis.Client) (userinfo.Store, error) {
	if err := db.Migrator().AutoMigrate(&userinfo.UserInfo{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) Insert(ctx context.Context, userID string, encryptionKeyChecksum string, preferedTimezone string) (userinfo.UserInfo, error) {
	userInfo := userinfo.UserInfo{
		ID:                    userID,
		PreferedTimeZone:      preferedTimezone,
		EncryptionKeyChecksum: encryptionKeyChecksum,
		StorageUsage:          0,
	}
	err := s.db.WithContext(ctx).Create(&userInfo).Error
	return userInfo, err
}

func (s *sqlStore) GetByID(ctx context.Context, userID string) (userInfo userinfo.UserInfo, err error) {
	err = s.db.Model(&userInfo).First(&userInfo, "id = ?", userID).Error
	return
}
