package uploadrequestsimpl

import (
	"context"

	uploadrequests "github.com/rishabhkailey/media-vault/internal/services/uploadRequests"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"gorm.io/gorm"
)

type FakeService struct {
	ExpectedError         error
	ExpectedUploadRequest storemodels.UploadRequestsModel
}

func NewFakeService() uploadrequests.Service {
	return &FakeService{}
}

var _ uploadrequests.Service = (*FakeService)(nil)

func (s *FakeService) WithTransaction(_ *gorm.DB) uploadrequests.Service {
	return s
}

func (s *FakeService) Create(ctx context.Context, cmd uploadrequests.CreateUploadRequestCommand) (storemodels.UploadRequestsModel, error) {
	return s.ExpectedUploadRequest, s.ExpectedError
}

func (s *FakeService) GetByID(ctx context.Context, query uploadrequests.GetByIDQuery) (storemodels.UploadRequestsModel, error) {
	return s.ExpectedUploadRequest, s.ExpectedError
}

func (s *FakeService) UpdateStatus(ctx context.Context, cmd uploadrequests.UpdateStatusCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteOne(ctx context.Context, cmd uploadrequests.DeleteOneCommand) error {
	return s.ExpectedError
}

func (s *FakeService) DeleteMany(ctx context.Context, cmd uploadrequests.DeleteManyCommand) error {
	return s.ExpectedError
}
