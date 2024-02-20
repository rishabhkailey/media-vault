package mediametadataimpl

import (
	"context"

	mediametadata "github.com/rishabhkailey/media-vault/internal/store/mediaMetadata"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

var _ mediametadata.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&storemodels.MediaMetadataModel{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) mediametadata.Store {
	return &sqlStore{
		db: tx,
	}
}

func (s *sqlStore) Insert(ctx context.Context, mediaMetadata *storemodels.MediaMetadataModel) (uint, error) {
	err := s.db.WithContext(ctx).Create(&mediaMetadata).Error
	return mediaMetadata.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, ID uint) error {
	return s.db.WithContext(ctx).Delete(&storemodels.MediaMetadataModel{
		Model: gorm.Model{
			ID: ID,
		},
	}).Error
}

func (s *sqlStore) DeleteMany(ctx context.Context, IDs []uint) error {
	return s.db.WithContext(ctx).Delete(&storemodels.MediaMetadataModel{}, IDs).Error
}

func (s *sqlStore) UpdateThumbnail(ctx context.Context, id uint, hasThumbnail bool, thumbnailAspectRatio float32) (err error) {
	err = s.db.Model(&storemodels.MediaMetadataModel{}).
		Where("id = ?", id).
		Select("thumbnail", "thumbnail_aspect_ratio").
		Updates(
			storemodels.MediaMetadataModel{
				MediaMetadata: storemodels.MediaMetadata{
					Thumbnail:            hasThumbnail,
					ThumbnailAspectRatio: thumbnailAspectRatio,
				},
			},
		).Error
	return
}
