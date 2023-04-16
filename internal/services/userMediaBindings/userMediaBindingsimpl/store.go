package usermediabindingsimpl

import (
	"context"
	"errors"
	"fmt"

	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"gorm.io/gorm"
)

type store interface {
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, *usermediabindings.Model) (uint, error)
	GetByMediaID(context.Context, uint) (usermediabindings.Model, error)
	CheckFileBelongsToUser(context.Context, usermediabindings.CheckFileBelongsToUserQuery) (bool, error)
	GetUserMedia(context.Context, usermediabindings.GetUserMediaQuery) ([]media.Model, error)
}

type sqlStore struct {
	db *gorm.DB
}

var _ store = (*sqlStore)(nil)

func newSqlStore(db *gorm.DB) (*sqlStore, error) {
	if err := db.Migrator().AutoMigrate(&usermediabindings.Model{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db: db,
	}, nil
}

func (s *sqlStore) Insert(ctx context.Context, userMediaBinding *usermediabindings.Model) (uint, error) {
	err := s.db.WithContext(ctx).Create(&userMediaBinding).Error
	return userMediaBinding.ID, err
}

func (s *sqlStore) GetByMediaID(ctx context.Context, mediaID uint) (userMediaBinding usermediabindings.Model, err error) {
	err = s.db.WithContext(ctx).First(&userMediaBinding, "media_id = ?", mediaID).Error
	return
}

func (s *sqlStore) CheckFileBelongsToUser(ctx context.Context, cmd usermediabindings.CheckFileBelongsToUserQuery) (ok bool, err error) {
	db := s.db.WithContext(ctx)
	getMediaByFileNameQuery := db.Model(&media.Model{}).Select("media_id").Where("file_name = ?", cmd.FileName)
	userMediaBinding := usermediabindings.Model{}
	err = db.Model(&usermediabindings.Model{}).Where("media_id = (?)", getMediaByFileNameQuery).First(&userMediaBinding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return
	}
	return userMediaBinding.UserID == cmd.UserID, nil
}

func (s *sqlStore) GetUserMedia(ctx context.Context, query usermediabindings.GetUserMediaQuery) (mediaList []media.Model, err error) {
	db := s.db.WithContext(ctx)
	mediaByUserIDQuery := db.Model(&usermediabindings.Model{}).Select("media_id").Where("user_id = ?", query.UserID)
	orderBy := fmt.Sprintf(`"Metadata"."%s"`, query.OrderBy)
	if query.Sort == usermediabindings.SORT_DESCENDING {
		orderBy = fmt.Sprintf("%s desc", orderBy)
	}
	err = db.Joins("Metadata").Model(&media.Model{}).Where("media.id IN (?)", mediaByUserIDQuery).Limit(query.Limit).Order(orderBy).Offset(query.Offset).Find(&mediaList).Error
	return

}
