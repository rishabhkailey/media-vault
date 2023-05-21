package uploadrequestsimpl

import (
	"context"

	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

type store interface {
	WithTransaction(*gorm.DB) store
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *uploadrequests.Model) (string, error)
	DeleteOne(context.Context, string) error
	DeleteMany(context.Context, []string) error
	GetByID(context.Context, string) (uploadrequests.Model, error)
	UpdateStatus(context.Context, uploadrequests.UpdateStatusCommand) error
}

type sqlStore struct {
	db *gorm.DB
}

var _ store = (*sqlStore)(nil)

func newSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&uploadrequests.Model{}); err != nil {
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

func (s *sqlStore) Insert(ctx context.Context, uploadRequest *uploadrequests.Model) (string, error) {
	err := s.db.WithContext(ctx).Create(&uploadRequest).Error
	return uploadRequest.ID, err
}

func (s *sqlStore) DeleteOne(ctx context.Context, id string) (err error) {
	return s.db.WithContext(ctx).Delete(&uploadrequests.Model{
		ID: id,
	}).Error
}

func (s *sqlStore) DeleteMany(ctx context.Context, ids []string) (err error) {
	return s.db.WithContext(ctx).Delete(&uploadrequests.Model{}, ids).Error
}

func (s *sqlStore) GetByID(ctx context.Context, id string) (uploadRequest uploadrequests.Model, err error) {
	err = s.db.WithContext(ctx).First(&uploadRequest, "id = ?", id).Error
	return
}

func (s *sqlStore) UpdateStatus(ctx context.Context, cmd uploadrequests.UpdateStatusCommand) (err error) {
	err = s.db.WithContext(ctx).Model(&uploadrequests.Model{ID: cmd.ID}).Updates(uploadrequests.Model{
		Status: cmd.Status,
	}).Error
	return
}
