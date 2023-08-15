package mediaimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	media "github.com/rishabhkailey/media-service/internal/store/media"
	"gorm.io/gorm"
)

// todo cache should have different store?
// maybe we can do something similar to media storage(cache wrapper for store)
type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

var _ media.Store = (*sqlStore)(nil)

// todo
// func newSqlStore(db *gorm.DB, cache *redis.Client) store {
// 	return &sqlStore{
// 		db:    db,
// 		cache: cache,
// 	}
// }

func NewSqlStoreWithMigrate(db *gorm.DB, cache *redis.Client) (media.Store, error) {
	if err := db.Migrator().AutoMigrate(&media.Media{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) media.Store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

func (s *sqlStore) Insert(ctx context.Context, media *media.Media) (uint, error) {
	err := s.db.WithContext(ctx).Create(&media).Error
	return media.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Delete(&media.Media{
		Model: gorm.Model{
			ID: id,
		},
	}).Error
	return err
}

func (s *sqlStore) DeleteMany(ctx context.Context, ids []uint) error {
	err := s.db.WithContext(ctx).Delete(&media.Media{}, ids).Error
	return err
}

func (s *sqlStore) GetByUploadRequestID(ctx context.Context, uploadRequestID string) (media media.Media, err error) {
	err = s.db.WithContext(ctx).First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetMediaWithMetadataByUploadRequestID(ctx context.Context, uploadRequestID string) (media media.Media, err error) {
	err = s.db.WithContext(ctx).Preload("Metadata").First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetByFileName(ctx context.Context, fileName string) (media media.Media, err error) {
	err = s.db.WithContext(ctx).First(&media, "file_name = ?", fileName).Error
	return
}

func (s *sqlStore) GetByUserIDOrderByDate(ctx context.Context,
	userID string,
	lastMediaID *uint,
	lastDate *time.Time,
	sort media.Sort,
	limit int,
) (mediaList []media.Media, err error) {

	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Model(&usermediabindings.Model{}).Select("media_id").Where("user_id = ?", userID)
	// table name = media, metadata alias = Metadata
	query := db.Joins("Metadata").Model(&media.Media{})
	if lastMediaID != nil && lastDate != nil {
		switch sort {
		case media.Ascending:
			{
				query = query.Where(`
				"media"."id" IN (?) 
				AND (
					(
						("Metadata"."date" = ?) AND ("media"."id" < ?)
					) OR (
						("Metadata"."date" > ?)
					)
					)`, mediaByUserIDQuery, lastDate, lastMediaID, lastDate)
			}
		default:
			{
				query = query.Where(`
				"media"."id" IN (?) 
				AND (
					(
						("Metadata"."date" = ?) AND ("media"."id" < ?)
					) OR (
						("Metadata"."date" < ?)
					)
					)`, mediaByUserIDQuery, lastDate, lastMediaID, lastDate)
			}
		}
	} else {
		query = query.Where(`
		"media"."id" IN (?) 
		`, mediaByUserIDQuery)
	}

	queryOrderBy := fmt.Sprintf(`"Metadata"."date" %s, "media"."id" desc`, sort)
	query = query.Order(queryOrderBy).Limit(limit)
	err = query.Find(&mediaList).Error
	return
}

func (s *sqlStore) GetByUserIDOrderByUploadDate(ctx context.Context,
	userID string,
	lastMediaID *uint,
	lastDate *time.Time,
	sort media.Sort,
	limit int,
) (mediaList []media.Media, err error) {

	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Model(&usermediabindings.Model{}).Select("media_id").Where("user_id = ?", userID)
	// table name = media, metadata alias = Metadata
	query := db.Joins("Metadata").Model(&media.Media{})
	if lastMediaID != nil && lastDate != nil {
		switch sort {
		case media.Ascending:
			{
				query = query.Where(`
				"media"."id" IN (?) 
				AND (
					(
						("Metadata"."created_at" = ?) AND ("media"."id" < ?)
					) OR (
						("Metadata"."created_at" > ?)
					)
					)`, mediaByUserIDQuery, lastDate, lastMediaID, lastDate)
			}
		default:
			{
				query = query.Where(`
				"media"."id" IN (?) 
				AND (
					(
						("Metadata"."created_at" = ?) AND ("media"."id" < ?)
					) OR (
						("Metadata"."created_at" < ?)
					)
					)`, mediaByUserIDQuery, lastDate, lastMediaID, lastDate)
			}
		}
	} else {
		query = query.Where(`
		"media"."id" IN (?) 
		`, mediaByUserIDQuery)
	}

	queryOrderBy := fmt.Sprintf(`"Metadata"."created_at" %s, "media"."id" desc`, sort)
	query = query.Order(queryOrderBy).Limit(limit)
	err = query.Find(&mediaList).Error
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
	media := media.Media{}
	err = db.Preload("Metadata").First(&media, "file_name = ?", fileName).Error
	if err == nil {
		mediaType = media.Metadata.Type
	}
	s.cache.Set(ctx, key, mediaType, 1*time.Hour)
	return
}

func (s *sqlStore) GetByMediaIDs(ctx context.Context, orderBy media.OrderBy, sort media.Sort, mediaIDs []uint) (mediaList []media.Media, err error) {
	order := fmt.Sprintf(`"Metadata"."%s" %s`, string(orderBy), string(sort))
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&media.Media{}).Where("media.id IN (?)", mediaIDs).Order(order).Find(&mediaList).Error
	return
}

func (s *sqlStore) GetByMediaID(ctx context.Context, mediaID uint) (m media.Media, err error) {
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&media.Media{}).Where("media.id = (?)", mediaID).First(&m).Error
	return
}
