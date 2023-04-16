package uploadrequestsimpl

import (
	"context"

	"github.com/google/uuid"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"gorm.io/gorm"
)

type Service struct {
	store store
}

var _ uploadrequests.Service = (*Service)(nil)

func NewService(db *gorm.DB) (uploadrequests.Service, error) {
	store, err := newSqlStore(db)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd uploadrequests.CreateUploadRequestCommand) (uploadrequests.Model, error) {
	UploadRequest := uploadrequests.Model{
		ID:     uuid.New().String(),
		UserID: cmd.UserID,
		Status: uploadrequests.IN_PROGRESS_UPLOAD_STATUS,
	}

	_, err := s.store.Insert(ctx, &UploadRequest)
	return UploadRequest, err
}

func (s *Service) GetByID(ctx context.Context, query uploadrequests.GetByIDQuery) (uploadrequests.Model, error) {
	return s.store.GetByID(ctx, query.ID)
}

func (s *Service) UpdateStatus(ctx context.Context, cmd uploadrequests.UpdateStatusCommand) error {
	return s.store.UpdateStatus(ctx, cmd)
}
