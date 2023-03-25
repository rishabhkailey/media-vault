package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ORDER_BY_UPLOAD_TIME = "created_at"
const ORDER_BY_MEDIA_CREATION_TIME = "time"

var SUPPORTED_ORDER_BY = []string{ORDER_BY_UPLOAD_TIME, ORDER_BY_MEDIA_CREATION_TIME}

const SORT_ASCENDING = "asc"
const SORT_DESCENDING = "desc"

// for now media will have 1-to-1 relationship with upload request and metadata
// todo make uploadrequest pointer so we can set the foreign key to null of media without upload request
type Media struct {
	gorm.Model
	FileName        string
	UploadRequestID string `gorm:"index"`
	UploadRequest   UploadRequest
	Metadata        MediaMetadata
	MetadataID      uint
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

func (model *MediaModel) Create(ctx context.Context, uploadRequestID string, metadataID uint) (*Media, error) {
	media := Media{
		FileName:        uuid.New().String(),
		UploadRequestID: uploadRequestID,
		MetadataID:      metadataID,
	}
	err := model.Db.WithContext(ctx).Create(&media).Error
	if err != nil {
		return nil, fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return &media, nil
}

func (model *MediaModel) FindByUploadRequest(ctx context.Context, uploadRequestID string) (media Media, err error) {
	err = model.Db.First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (model *MediaModel) GetUserMediaList(ctx context.Context, userID string, orderBy string, sort string, offset, limit int) (mediaList []Media, err error) {
	db := model.Db
	mediaIdQuery := db.Model(&UserMediaBinding{}).Select("media_id").Where("user_id = ?", userID)
	if sort == SORT_DESCENDING {
		orderBy = fmt.Sprintf("%s desc", orderBy)
	}
	// todo check preload vs join
	err = db.Preload("Metadata").Model(&Media{}).Where("id IN (?)", mediaIdQuery).Limit(limit).Order(orderBy).Offset(offset).Find(&mediaList).Error
	return
}

// todo redis for user media count update on each upload and delete operation
// 1 more counter for failed uploads if user wants to see

// todo check gorm scopes https://gorm.io/docs/advanced_query.html#Scopes
// todo check pagination https://gorm.io/docs/scopes.html#Pagination
