package usermediabindingsimpl

import (
	"context"
	"errors"
	"fmt"

	// usermediabindings "github.com/rishabhkailey/media-vault/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-vault/internal/constants"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	usermediabindings "github.com/rishabhkailey/media-vault/internal/store/userMediaBindings"
	"github.com/rishabhkailey/media-vault/internal/utils"
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

var _ usermediabindings.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&storemodels.UserMediaBindingsModel{}); err != nil {
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

func (s *sqlStore) Insert(ctx context.Context, userMediaBinding *storemodels.UserMediaBindingsModel) (uint, error) {
	err := s.db.WithContext(ctx).Create(&userMediaBinding).Error
	return userMediaBinding.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, userID string, mediaID uint) error {
	return s.db.WithContext(ctx).Delete(&storemodels.UserMediaBindingsModel{
		UserID:  userID,
		MediaID: mediaID,
	}).Error
}

// func (s *sqlStore) DeleteMany(ctx context.Context, userID string, mediaIDs []uint) error {
// 	return s.db.WithContext(ctx).Where("media_id IN ?", mediaIDs).Delete(&usermediabindings.Model{
// 		UserID: userID,
// 	}).Error
// }

func (s *sqlStore) GetByMediaID(ctx context.Context, mediaID uint) (userMediaBinding storemodels.UserMediaBindingsModel, err error) {
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
	userMediaBinding := storemodels.UserMediaBindingsModel{}
	err = db.Model(&storemodels.UserMediaBindingsModel{}).
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

func (s *sqlStore) CheckMultipleMediaBelongsToUser(ctx context.Context, userID string, mediaIDs []uint) (bool, error) {
	var userOwnedMediaIds []uint
	err := s.db.WithContext(ctx).Model(&storemodels.UserMediaBindingsModel{}).
		Select("media_id").
		Where("user_id = ? AND media_id IN (?)", userID, mediaIDs).
		Find(&userOwnedMediaIds).Error
	if err != nil {
		return false, fmt.Errorf("[usermediabindingsimpl.CheckMultipleMediaBelongsToUser] failed to check user access: %w", err)
	}
	if !utils.ContainsSlice(userOwnedMediaIds, mediaIDs) {
		return false, nil
	}
	return true, nil
}
