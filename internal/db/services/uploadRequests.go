package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadRequest struct {
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"index:,composite:user_id_status"`
	Status Status `gorm:"index:,composite:user_id_status"`
}

type UploadRequestModel struct {
	Db *gorm.DB
}

type Status string

const COMPLETED_UPLOAD_STATUS Status = "completed"
const FAILED_UPLOAD_STATUS Status = "failed"
const IN_PROGRESS_UPLOAD_STATUS Status = "inProgress"

func NewUploadRequestModel(db *gorm.DB) (*UploadRequestModel, error) {
	err := db.Migrator().AutoMigrate(&UploadRequest{})
	if err != nil {
		return nil, err
	}
	return &UploadRequestModel{
		Db: db,
	}, nil
}

func (model *UploadRequestModel) FindByID(ctx context.Context, id string) (UploadRequest UploadRequest, err error) {
	err = model.Db.First(&UploadRequest, "id = ?", id).Error
	if err != nil {
		return UploadRequest, fmt.Errorf("error getting UploadRequest details from DB: %w", err)
	}
	return UploadRequest, err
}

func (model *UploadRequestModel) FindByEmail(ctx context.Context, email string) (UploadRequest UploadRequest, err error) {
	err = model.Db.First(&UploadRequest, "email = ?", email).Error
	if err != nil {
		return UploadRequest, fmt.Errorf("error getting UploadRequest details from DB: %w", err)
	}
	return UploadRequest, err
}

func (model *UploadRequestModel) Create(ctx context.Context, userID string) (*UploadRequest, error) {
	UploadRequest := &UploadRequest{
		ID:     uuid.New().String(),
		UserID: userID,
		Status: IN_PROGRESS_UPLOAD_STATUS,
	}
	err := model.Db.WithContext(ctx).Create(UploadRequest).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return UploadRequest, nil
}

func (model *UploadRequestModel) UpdateStatus(ctx context.Context, uploadRequestId string, status Status) error {
	err := model.Db.WithContext(ctx).Model(&UploadRequest{ID: uploadRequestId}).Updates(UploadRequest{
		Status: status,
	}).Error
	if err != nil {
		return fmt.Errorf("[Create]: update failed: %w", err)
	}
	return nil
}
