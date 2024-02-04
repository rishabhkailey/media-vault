package uploadrequestsimpl

import (
	"context"

	"github.com/google/uuid"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"github.com/rishabhkailey/media-service/internal/store"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
)

type Service struct {
	store store.Store
}

var _ uploadrequests.Service = (*Service)(nil)

func NewService(store store.Store) (uploadrequests.Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd uploadrequests.CreateUploadRequestCommand) (storemodels.UploadRequestsModel, error) {
	UploadRequest := storemodels.UploadRequestsModel{
		ID:     uuid.New().String(),
		UserID: cmd.UserID,
		Status: string(uploadrequests.IN_PROGRESS_UPLOAD_STATUS),
	}

	_, err := s.store.UploadRequests.Insert(ctx, &UploadRequest)
	return UploadRequest, err
}

func (s *Service) GetByID(ctx context.Context, query uploadrequests.GetByIDQuery) (storemodels.UploadRequestsModel, error) {
	return s.store.UploadRequests.GetByID(ctx, query.ID)
}

func (s *Service) DeleteOne(ctx context.Context, cmd uploadrequests.DeleteOneCommand) error {
	return s.store.UploadRequests.DeleteOne(ctx, cmd.ID)
}

func (s *Service) DeleteMany(ctx context.Context, cmd uploadrequests.DeleteManyCommand) error {
	return s.store.UploadRequests.DeleteMany(ctx, cmd.IDs)
}

func (s *Service) UpdateStatus(ctx context.Context, cmd uploadrequests.UpdateStatusCommand) error {
	return s.store.UploadRequests.UpdateStatus(ctx, cmd.ID, string(cmd.Status))
}
