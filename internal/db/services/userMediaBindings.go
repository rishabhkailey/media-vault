package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// todo user subject as userid on the client so even if user change email or user name the user should not loose its data
type UserMediaBinding struct {
	gorm.Model
	UserID  string `gorm:"index:,unique,composite:user_id_media_id"`
	MediaID string `gorm:"index:,unique,composite:user_id_media_id"`
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

func (model *UserMediaBindingModel) Create(ctx context.Context, userID string, media Media) (*UserMediaBinding, error) {
	userMediaBinding := &UserMediaBinding{
		UserID: userID,
		Media:  media,
	}
	err := model.Db.WithContext(ctx).Create(userMediaBinding).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return userMediaBinding, nil
}
