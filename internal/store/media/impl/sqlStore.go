package mediaimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rishabhkailey/media-vault/internal/constants"
	media "github.com/rishabhkailey/media-vault/internal/store/media"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
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
	if err := db.Migrator().AutoMigrate(&storemodels.MediaModel{}); err != nil {
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

func (s *sqlStore) Insert(ctx context.Context, media *storemodels.MediaModel) (uint, error) {
	err := s.db.WithContext(ctx).Create(&media).Error
	return media.ID, err
}

func (s *sqlStore) CascadeDeleteOne(ctx context.Context, mediaID uint, userID string, mediaMetadataId uint) error {
	tx := s.db.WithContext(ctx).Begin()
	err := tx.
		Where("user_id = ? AND media_id = ?", userID, mediaID).
		Delete(&storemodels.UserMediaBindingsModel{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[mediaimpl.DeleteOne] usermedia binding deletion failed: %w", err)
	}

	err = tx.Where("media_id = ?", mediaID).
		Delete(&storemodels.AlbumMediaBindingsModel{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[mediaimpl.DeleteOne] albummedia binding deletion failed: %w", err)
	}

	err = tx.Where("id = ?", mediaID).
		Delete(&storemodels.MediaModel{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[mediaimpl.DeleteOne] media deletion failed: %w", err)
	}

	err = tx.Where("id = ?", mediaMetadataId).
		Delete(&storemodels.MediaMetadataModel{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[mediaimpl.DeleteOne] MediaMetadataModel deletion failed: %w", err)
	}

	tx.Commit()
	return err
}

func (s *sqlStore) DeleteMany(ctx context.Context, ids []uint) error {
	err := s.db.WithContext(ctx).Delete(&storemodels.MediaModel{}, ids).Error
	return err
}

func (s *sqlStore) GetByUploadRequestID(ctx context.Context, uploadRequestID string) (media storemodels.MediaModel, err error) {
	err = s.db.WithContext(ctx).First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetMediaWithMetadataByUploadRequestID(ctx context.Context, uploadRequestID string) (media storemodels.MediaModel, err error) {
	err = s.db.WithContext(ctx).Preload("Metadata").First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetByFileName(ctx context.Context, fileName string) (media storemodels.MediaModel, err error) {
	err = s.db.WithContext(ctx).First(&media, "file_name = ?", fileName).Error
	return
}

func (s *sqlStore) GetByUserIDOrderByDate(ctx context.Context,
	userID string,
	lastMediaID *uint,
	lastDate *time.Time,
	sort media.Sort,
	limit int,
) (mediaList []storemodels.MediaModel, err error) {

	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Table(constants.USER_MEDIA_BINDINGS_TABLE).Select("media_id").Where("user_id = ?", userID)
	// table name = media, metadata alias = Metadata
	query := db.Joins("Metadata").Model(&storemodels.MediaModel{})
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
) (mediaList []storemodels.MediaModel, err error) {

	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Table(constants.USER_MEDIA_BINDINGS_TABLE).Select("media_id").Where("user_id = ?", userID)
	// table name = media, metadata alias = Metadata
	query := db.Joins("Metadata").Model(&storemodels.MediaModel{})
	if lastMediaID != nil && lastDate != nil {
		switch sort {
		case media.Ascending:
			{
				// media whose created date > last media date (sort by created_at asc)
				// media whose created date = last media date and media id < last media id (media id is sorted by desc order)
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
				// media whose created date < last media date (sort by created_at desc)
				// media whose created date = last media date and media id < last media id (media id is sorted by desc order)
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
	media := storemodels.MediaModel{}
	err = db.Preload("Metadata").First(&media, "file_name = ?", fileName).Error
	if err == nil {
		mediaType = media.Metadata.Type
	}
	s.cache.Set(ctx, key, mediaType, 1*time.Hour)
	return
}

func (s *sqlStore) GetByMediaIDsWithSort(ctx context.Context, orderBy media.OrderBy, sort media.Sort, mediaIDs []uint) (mediaList []storemodels.MediaModel, err error) {
	order := fmt.Sprintf(`"Metadata"."%s" %s`, string(orderBy), string(sort))
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&storemodels.MediaModel{}).Where("media.id IN (?)", mediaIDs).Order(order).Find(&mediaList).Error
	return
}

func (s *sqlStore) GetByMediaID(ctx context.Context, mediaID uint) (m storemodels.MediaModel, err error) {
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&storemodels.MediaModel{}).Where("media.id = (?)", mediaID).First(&m).Error
	return
}

func (s *sqlStore) GetByMediaIDs(ctx context.Context, mediaIDs []uint) (mediaList []storemodels.MediaModel, err error) {
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&storemodels.MediaModel{}).Where("media.id IN (?)", mediaIDs).Find(&mediaList).Error
	return
}

func (s *sqlStore) CascadeDeleteMany(ctx context.Context, userID string, mediaIDs []uint) (
	deletedUserMediaBindings []storemodels.UserMediaBindingsModel,
	deletedAlbumMediaBindings []storemodels.AlbumMediaBindingsModel,
	deletedMedia []storemodels.MediaModel,
	deletedMediaMetadata []storemodels.MediaMetadataModel,
	err error,
) {

	tx := s.db.WithContext(ctx).Begin()
	err = tx.
		Where("user_id = ? AND media_id IN (?)", userID, mediaIDs).
		Delete(&deletedUserMediaBindings).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("[mediaimpl.DeleteOne] usermedia binding deletion failed: %w", err)
		return
	}

	err = tx.Where("media_id IN (?)", mediaIDs).
		Delete(&deletedAlbumMediaBindings).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("[mediaimpl.DeleteOne] albummedia binding deletion failed: %w", err)
		return
	}

	err = tx.Where("id IN (?)", mediaIDs).
		Delete(&deletedMedia).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("[mediaimpl.DeleteOne] media deletion failed: %w", err)
		return
	}

	var mediaMetadataIDs []uint
	for _, media := range deletedMedia {
		mediaMetadataIDs = append(mediaMetadataIDs, media.MetadataID)
	}
	err = tx.Where("id IN (?)", mediaMetadataIDs).
		Delete(&deletedMediaMetadata).Error
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("[mediaimpl.DeleteOne] MediaMetadataModel deletion failed: %w", err)
		return
	}

	tx.Commit()
	return
}
