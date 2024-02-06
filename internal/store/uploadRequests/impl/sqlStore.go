package uploadrequestsimpl

import (
	"context"

	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	uploadrequests "github.com/rishabhkailey/media-service/internal/store/uploadRequests"
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

var _ uploadrequests.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&storemodels.UploadRequestsModel{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) uploadrequests.Store {
	return &sqlStore{
		db: tx,
	}
}

func (s *sqlStore) Insert(ctx context.Context, uploadRequest *storemodels.UploadRequestsModel) (string, error) {
	err := s.db.WithContext(ctx).Create(&uploadRequest).Error
	return uploadRequest.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, id string) (err error) {
	return s.db.WithContext(ctx).Delete(&storemodels.UploadRequestsModel{
		ID: id,
	}).Error
}

func (s *sqlStore) DeleteMany(ctx context.Context, ids []string) (err error) {
	return s.db.WithContext(ctx).Delete(&storemodels.UploadRequestsModel{}, ids).Error
}

func (s *sqlStore) GetByID(ctx context.Context, id string) (uploadRequest storemodels.UploadRequestsModel, err error) {
	err = s.db.WithContext(ctx).First(&uploadRequest, "id = ?", id).Error
	return
}

func (s *sqlStore) UpdateStatus(ctx context.Context, id string, status string) (err error) {
	err = s.db.WithContext(ctx).Model(&storemodels.UploadRequestsModel{ID: id}).
		Updates(storemodels.UploadRequestsModel{
			Status: status,
		}).Error
	return
}
