package albumimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgconn"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	"github.com/rishabhkailey/media-service/internal/store/album"
	"github.com/rishabhkailey/media-service/internal/store/media"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) album.Store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

var _ album.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB, cache *redis.Client) (album.Store, error) {
	if err := db.Migrator().AutoMigrate(&album.Album{}); err != nil {
		return nil, err
	}
	if err := db.Migrator().AutoMigrate(&album.AlbumMediaBindings{}); err != nil {
		return nil, err
	}
	if err := db.Migrator().AutoMigrate(&album.UserAlbumBindings{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) InsertAlbum(ctx context.Context, albumName string, thumbnailUrl string) (album.Album, error) {
	album := album.Album{
		Name:         albumName,
		ThumbnailUrl: thumbnailUrl,
	}
	err := s.db.WithContext(ctx).Create(&album).Error
	return album, err
}

func (s *sqlStore) InsertUserAlbumBindings(ctx context.Context, userID string, albumID uint) (uint, error) {
	album := album.UserAlbumBindings{
		UserID:  userID,
		AlbumID: albumID,
	}
	err := s.db.WithContext(ctx).Create(&album).Error
	return album.ID, err
}

func (s *sqlStore) GetByUserId(ctx context.Context, userID string, orderBy string, sort string, limit int, offset int) (albums []album.Album, err error) {
	db := s.db.WithContext(ctx)
	// todo join?
	// cost of join
	// cost of join > cost of fetching all media_ids of the album and then filter
	albumByUserIDQuery := db.Model(&album.UserAlbumBindings{}).Select("album_id").Where("user_id = ?", userID)
	err = db.Model(&album.Album{}).
		Where("id IN (?)", albumByUserIDQuery).
		Limit(limit).
		Order(fmt.Sprintf("%s %s", orderBy, sort)).
		Offset(offset).
		Find(&albums).Error
	return
}

func (s *sqlStore) GetByID(ctx context.Context, albumID uint) (result album.Album, err error) {
	err = s.db.WithContext(ctx).Model(&album.Album{}).First(&result, albumID).Error
	return
}

func (s *sqlStore) GetMediaByAlbumId(ctx context.Context, albumID uint, orderBy string, sort string, limit int, offset int) (mediaList []media.Media, err error) {
	db := s.db.WithContext(ctx)
	mediaIDsByAlbumQuery := db.Model(&album.AlbumMediaBindings{}).Select("media_id").Where("album_id = ?", albumID)
	err = db.Joins("Metadata").Model(&media.Media{}).
		Where("media.id IN (?)", mediaIDsByAlbumQuery).
		Limit(limit).
		Order(fmt.Sprintf("%s %s", orderBy, sort)).
		Offset(offset).
		Find(&mediaList).Error
	return
}

func (s *sqlStore) CheckAlbumBelongsToUser(ctx context.Context, userID string, albumID uint) (ok bool, err error) {
	var userAlbumBinding album.UserAlbumBindings
	err = s.db.WithContext(ctx).First(&userAlbumBinding, "album_id = ?", albumID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return userAlbumBinding.UserID == userID, nil
}

func (s *sqlStore) AddMediaInAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (newMediaIDs []uint, err error) {
	alreadyExist, err := s.FilterMediaIDsByAlbumID(ctx, albumID, mediaIDs)
	if err != nil {
		return
	}
	newMediaIDs = make([]uint, 0)
	var albumMediaBindings []album.AlbumMediaBindings
	for _, mediaID := range mediaIDs {
		if utils.Contains(alreadyExist, mediaID) {
			continue
		}
		newMediaIDs = append(newMediaIDs, mediaID)
		albumMediaBindings = append(albumMediaBindings, album.AlbumMediaBindings{
			AlbumID: albumID,
			MediaID: mediaID,
		})
	}
	if len(albumMediaBindings) == 0 {
		return
	}
	err = s.db.WithContext(ctx).CreateInBatches(albumMediaBindings, 100).Error
	if isUniqueConstraintError(err) {
		logrus.Warnf("[AddMediaInAlbum] duplicate constraint ignoring some rows: %v", err)
		err = nil
	}
	if err != nil {
		return
	}
	if _, err := s.UpdateAlbumMediaCount(ctx, albumID, len(newMediaIDs)); err != nil {
		logrus.Warnf("[AddMediaInAlbum] album updatedAt failed: %v", err)
	}
	return
}

func (s *sqlStore) RemoveMediaFromAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (removedMediaIDs []uint, err error) {
	mediaIDs, err = s.FilterMediaIDsByAlbumID(ctx, albumID, mediaIDs)
	if err != nil {
		return
	}
	db := s.db.WithContext(ctx).Begin()
	err = db.Unscoped().Delete(&album.AlbumMediaBindings{}, "album_id = ? AND media_id IN ?", albumID, mediaIDs).Error
	if err != nil {
		db.Rollback()
		return
	}
	db.Commit()
	removedMediaIDs = mediaIDs
	if _, err := s.UpdateAlbumMediaCount(ctx, albumID, len(removedMediaIDs)*-1); err != nil {
		logrus.Warnf("[AddMediaInAlbum] album updatedAt failed: %v", err)
	}
	return
}

func (s *sqlStore) FilterMediaIDsByAlbumID(ctx context.Context, albumID uint, mediaIDs []uint) (result []uint, err error) {
	err = s.db.Model(&album.AlbumMediaBindings{}).Where("media_id IN ? and album_id = ?", mediaIDs, albumID).Pluck("media_id", &result).Error
	return
}

func (s *sqlStore) DeleteAlbum(ctx context.Context, albumID uint) error {
	tx := s.db.Begin().WithContext(ctx)

	if err := tx.Unscoped().Delete(&album.UserAlbumBindings{}, "album_id = ?", albumID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&album.AlbumMediaBindings{}, "album_id = ?", albumID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&album.Album{
		Model: gorm.Model{
			ID: albumID,
		},
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *sqlStore) UpdateAlbum(ctx context.Context, albumID uint, name *string, thumbnailUrl *string) (updatedAlbum album.Album, err error) {
	err = s.db.WithContext(ctx).First(&updatedAlbum, albumID).Error
	if err != nil {
		return
	}
	if name != nil {
		updatedAlbum.Name = *name
	}
	if thumbnailUrl != nil {
		updatedAlbum.ThumbnailUrl = *thumbnailUrl
	}
	err = s.db.WithContext(ctx).Save(&updatedAlbum).Error
	return
}

func (s *sqlStore) UpdateAlbumMediaCount(ctx context.Context, albumID uint, change int) (updatedAlbum album.Album, err error) {
	// album.ID = albumID
	// album.UpdatedAt = updatedAt

	err = s.db.WithContext(ctx).First(&updatedAlbum, albumID).Error
	if err != nil {
		return
	}
	updatedAlbum.MediaCount += change
	if updatedAlbum.MediaCount < 0 {
		updatedAlbum.MediaCount = 0
	}
	err = s.db.WithContext(ctx).Save(&updatedAlbum).Error
	return
}

func (s *sqlStore) UpdateThumbnail(ctx context.Context, mediaID uint, thumbnail bool, thumbnailAspectRatio float32) error {
	mediaMetadataIdSubQuery := s.db.Model(&media.Media{}).Select("metadata_id").Where("id = ?", mediaID).Limit(1)
	return s.db.Model(&mediametadata.Model{}).Where("id = ?", mediaMetadataIdSubQuery).Select("thumbnail", "thumbnail_aspect_ratio").Updates(
		mediametadata.Model{
			Metadata: mediametadata.Metadata{
				Thumbnail:            thumbnail,
				ThumbnailAspectRatio: thumbnailAspectRatio,
			},
		},
	).Error
}

func isUniqueConstraintError(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		if pgError.Code == "23505" {
			return true
		}
	}
	return false
}
