package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	FileName        string
	UploadRequestID string
	UploadRequest   UploadRequest
}

type MediaModel struct {
	Db *gorm.DB
}

func NewMediaModel(db *gorm.DB) (*MediaModel, error) {
	err := db.Migrator().AutoMigrate(&Media{})
	if err != nil {
		return nil, err
	}
	return &MediaModel{
		Db: db,
	}, nil
}

func (model *MediaModel) Create(ctx context.Context, fileName string, uploadRequest UploadRequest) (*Media, error) {
	media := Media{
		FileName:      fileName,
		UploadRequest: uploadRequest,
	}
	err := model.Db.WithContext(ctx).Create(&media).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return &media, nil
}
