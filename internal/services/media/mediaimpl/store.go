package mediaimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"gorm.io/gorm"
)

type store interface {
	WithTransaction(*gorm.DB) store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *media.Model) (uint, error)
	DeleteOne(context.Context, uint) error
	DeleteMany(context.Context, []uint) error
	GetByUploadRequestID(context.Context, string) (media.Model, error)
	GetByFileName(context.Context, string) (media.Model, error)
	GetMediaWithMetadataByUploadRequestID(context.Context, string) (media.Model, error)
	GetByUserID(context.Context, media.GetByUserIDQuery) ([]media.Model, error)
	GetByMediaID(context.Context, media.GetByMediaIDQuery) (media.Model, error)
	GetByMediaIDs(context.Context, media.GetByMediaIDsQuery) ([]media.Model, error)
	GetTypeByFileName(context.Context, string) (string, error)
}

// todo cache should have different store?
// maybe we can do something similar to media storage(cache wrapper for store)
type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

var _ store = (*sqlStore)(nil)

// todo
// func newSqlStore(db *gorm.DB, cache *redis.Client) store {
// 	return &sqlStore{
// 		db:    db,
// 		cache: cache,
// 	}
// }

func newSqlStoreWithMigrate(db *gorm.DB, cache *redis.Client) (store, error) {
	if err := db.Migrator().AutoMigrate(&media.Model{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

func (s *sqlStore) Insert(ctx context.Context, media *media.Model) (uint, error) {
	err := s.db.WithContext(ctx).Create(&media).Error
	return media.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Delete(&media.Model{
		Model: gorm.Model{
			ID: id,
		},
	}).Error
	return err
}

func (s *sqlStore) DeleteMany(ctx context.Context, ids []uint) error {
	err := s.db.WithContext(ctx).Delete(&media.Model{}, ids).Error
	return err
}

func (s *sqlStore) GetByUploadRequestID(ctx context.Context, uploadRequestID string) (media media.Model, err error) {
	err = s.db.WithContext(ctx).First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetMediaWithMetadataByUploadRequestID(ctx context.Context, uploadRequestID string) (media media.Model, err error) {
	err = s.db.WithContext(ctx).Preload("Metadata").First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetByFileName(ctx context.Context, fileName string) (media media.Model, err error) {
	err = s.db.WithContext(ctx).First(&media, "file_name = ?", fileName).Error
	return
}

func (s *sqlStore) GetByUserID(ctx context.Context, query media.GetByUserIDQuery) (mediaList []media.Model, err error) {
	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Model(&usermediabindings.Model{}).Select("media_id").Where("user_id = ?", query.UserID)
	orderBy := fmt.Sprintf(`"Metadata"."%s" %s`, media.OrderAttributesMapping[query.OrderBy], media.SortKeywordMapping[query.Sort])
	limit := int(query.PerPage)
	offset := int((query.Page - 1) * query.PerPage)
	err = db.Joins("Metadata").Model(&media.Model{}).Where("media.id IN (?)", mediaByUserIDQuery).Limit(limit).Order(orderBy).Offset(offset).Find(&mediaList).Error
	return
}

// todo instead just use minio file type ? but other storage types may not have this option
func (s *sqlStore) GetTypeByFileName(ctx context.Context, fileName string) (mediaType string, err error) {
	key := fmt.Sprintf("mediaType:%s", fileName)
	mediaType, err = s.cache.Get(ctx, key).Result()
	if err == nil {
		return
	}
	db := s.db.WithContext(ctx)
	media := media.Model{}
	err = db.Preload("Metadata").First(&media, "file_name = ?", fileName).Error
	if err == nil {
		mediaType = media.Metadata.Type
	}
	s.cache.Set(ctx, key, mediaType, 1*time.Hour)
	return
}

func (s *sqlStore) GetByMediaIDs(ctx context.Context, query media.GetByMediaIDsQuery) (mediaList []media.Model, err error) {
	orderBy := fmt.Sprintf(`"Metadata"."%s" %s`, media.OrderAttributesMapping[query.OrderBy], media.SortKeywordMapping[query.Sort])
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&media.Model{}).Where("media.id IN (?)", query.MediaIDs).Order(orderBy).Find(&mediaList).Error
	return
}

func (s *sqlStore) GetByMediaID(ctx context.Context, query media.GetByMediaIDQuery) (m media.Model, err error) {
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&media.Model{}).Where("media.id = (?)", query.MediaID).First(&m).Error
	return
}
