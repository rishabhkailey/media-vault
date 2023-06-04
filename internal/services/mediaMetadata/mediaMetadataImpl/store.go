package mediametadataimpl

import (
	"context"

	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	"gorm.io/gorm"
)

type store interface {
	WithTransaction(*gorm.DB) store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *mediametadata.Model) (uint, error)
	DeleteOne(context.Context, uint) error
	DeleteMany(context.Context, []uint) error
	UpdateThumbnail(context.Context, mediametadata.UpdateThumbnailCommand) error
}

type sqlStore struct {
	db *gorm.DB
}

var _ store = (*sqlStore)(nil)

func newSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&mediametadata.Model{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) store {
	return &sqlStore{
		db: tx,
	}
}

func (s *sqlStore) Insert(ctx context.Context, mediaMetadata *mediametadata.Model) (uint, error) {
	err := s.db.WithContext(ctx).Create(&mediaMetadata).Error
	return mediaMetadata.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, ID uint) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&mediametadata.Model{
		Model: gorm.Model{
			ID: ID,
		},
	}).Error
}

func (s *sqlStore) DeleteMany(ctx context.Context, IDs []uint) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&mediametadata.Model{}, IDs).Error
}

func (s *sqlStore) UpdateThumbnail(ctx context.Context, cmd mediametadata.UpdateThumbnailCommand) (err error) {
	err = s.db.Model(&mediametadata.Model{}).Where("id = ?", cmd.ID).Update("thumbnail", cmd.Thumbnail).Error
	return
}
