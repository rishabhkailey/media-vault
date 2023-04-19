package mediaimpl

import (
	"context"
	"fmt"

	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"gorm.io/gorm"
)

type store interface {
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *media.Model) (uint, error)
	GetByUploadRequestID(context.Context, string) (media.Model, error)
	GetByFileName(context.Context, string) (media.Model, error)
	GetByUserID(context.Context, media.GetByUserIDQuery) ([]media.Model, error)
	GetByMediaIDs(context.Context, []uint) ([]media.Model, error)
	GetTypeByFileName(context.Context, string) (string, error)
}

type sqlStore struct {
	db *gorm.DB
}

var _ store = (*sqlStore)(nil)

func newSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&media.Model{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) Insert(ctx context.Context, media *media.Model) (uint, error) {
	err := s.db.WithContext(ctx).Create(&media).Error
	return media.ID, err
}

func (s *sqlStore) GetByUploadRequestID(ctx context.Context, uploadRequestID string) (media media.Model, err error) {
	err = s.db.WithContext(ctx).First(&media, "upload_request_id = ?", uploadRequestID).Error
	return
}

func (s *sqlStore) GetByFileName(ctx context.Context, fileName string) (media media.Model, err error) {
	err = s.db.WithContext(ctx).First(&media, "file_name = ?", fileName).Error
	return
}

func (s *sqlStore) GetByUserID(ctx context.Context, query media.GetByUserIDQuery) (mediaList []media.Model, err error) {
	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Model(&usermediabindings.Model{}).Select("media_id").Where("user_id = ?", query.UserID)
	orderBy := fmt.Sprintf(`"Metadata"."%s"`, query.OrderBy)
	if query.Sort == media.SORT_DESCENDING {
		orderBy = fmt.Sprintf("%s desc", orderBy)
	}
	err = db.Joins("Metadata").Model(&media.Model{}).Where("media.id IN (?)", mediaByUserIDQuery).Limit(query.Limit).Order(orderBy).Offset(query.Offset).Find(&mediaList).Error
	return
}

func (s *sqlStore) GetTypeByFileName(ctx context.Context, fileName string) (mediaType string, err error) {
	db := s.db.WithContext(ctx)
	media := media.Model{}
	err = db.Preload("Metadata").First(&media, "file_name = ?", fileName).Error
	if err == nil {
		mediaType = media.Metadata.Type
	}
	return
}

func (s *sqlStore) GetByMediaIDs(ctx context.Context, mediaIDs []uint) (mediaList []media.Model, err error) {
	err = s.db.WithContext(ctx).Joins("Metadata").Model(&media.Model{}).Where("media.id IN (?)", mediaIDs).Find(&mediaList).Error
	return
}
