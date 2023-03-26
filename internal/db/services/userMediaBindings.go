package services

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// todo user subject as userid on the client so even if user change email or user name the user should not loose its data
type UserMediaBinding struct {
	gorm.Model
	UserID  string `gorm:"index:,unique,composite:user_id_media_id"`
	MediaID uint   `gorm:"index:,unique,composite:user_id_media_id"`
	Media   Media  `gorm:"foreignKey:MediaID"`
}

type UserMediaBindingModel struct {
	Db *gorm.DB
}

func NewUserMediaBinding(db *gorm.DB) (*UserMediaBindingModel, error) {
	err := db.Migrator().AutoMigrate(&UserMediaBinding{})
	if err != nil {
		return nil, err
	}
	return &UserMediaBindingModel{
		Db: db,
	}, nil
}

func (model *UserMediaBindingModel) Create(ctx context.Context, userID string, mediaID uint) (*UserMediaBinding, error) {
	userMediaBinding := &UserMediaBinding{
		UserID:  userID,
		MediaID: mediaID,
	}
	err := model.Db.WithContext(ctx).Create(userMediaBinding).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return userMediaBinding, nil
}

func (model *UserMediaBindingModel) FindByMediaID(ctx context.Context, mediaID string) (userMediaBinding UserMediaBinding, err error) {
	err = model.Db.First(&userMediaBinding, "media_id = ?", mediaID).Error
	return
}

func (model *UserMediaBindingModel) CheckFileBelongsToUser(ctx context.Context, userID, fileName string) (ok bool, err error) {
	db := model.Db
	getMediaByFileNameQuery := db.Model(&Media{}).Select("media_id").Where("file_name = ?", fileName)
	userMediaBinding := UserMediaBinding{}
	err = db.Model(&UserMediaBinding{}).Where("media_id = (?)", getMediaByFileNameQuery).First(&userMediaBinding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return
	}
	return userMediaBinding.UserID == userID, nil
}
