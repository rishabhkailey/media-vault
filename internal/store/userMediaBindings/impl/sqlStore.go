package usermediabindingsimpl

import (
	"context"
	"errors"

	// usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-service/internal/constants"
	usermediabindings "github.com/rishabhkailey/media-service/internal/store/userMediaBindings"
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

var _ usermediabindings.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&usermediabindings.Model{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) usermediabindings.Store {
	return &sqlStore{
		db: tx,
	}
}

func (s *sqlStore) Insert(ctx context.Context, userMediaBinding *usermediabindings.Model) (uint, error) {
	err := s.db.WithContext(ctx).Create(&userMediaBinding).Error
	return userMediaBinding.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, userID string, mediaID uint) error {
	return s.db.WithContext(ctx).Delete(&usermediabindings.Model{
		UserID:  userID,
		MediaID: mediaID,
	}).Error
}

// func (s *sqlStore) DeleteMany(ctx context.Context, userID string, mediaIDs []uint) error {
// 	return s.db.WithContext(ctx).Where("media_id IN ?", mediaIDs).Delete(&usermediabindings.Model{
// 		UserID: userID,
// 	}).Error
// }

func (s *sqlStore) GetByMediaID(ctx context.Context, mediaID uint) (userMediaBinding usermediabindings.Model, err error) {
	err = s.db.WithContext(ctx).First(&userMediaBinding, "media_id = ?", mediaID).Error
	return
}

// todo this logic should be in service not in store
// it should be like get by fileName
// and service should check the userID
func (s *sqlStore) CheckFileBelongsToUser(ctx context.Context, userID, fileName string) (ok bool, err error) {
	db := s.db.WithContext(ctx)
	getMediaByFileNameQuery := db.Table(constants.MEDIA_TABLE).
		Select("id").
		Where("file_name = ?", fileName).
		Limit(1)
	userMediaBinding := usermediabindings.Model{}
	err = db.Model(&usermediabindings.Model{}).
		Where("media_id = (?)", getMediaByFileNameQuery).
		First(&userMediaBinding).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return
	}
	return userMediaBinding.UserID == userID, nil
}
