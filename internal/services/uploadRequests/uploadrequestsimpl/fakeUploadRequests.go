package uploadrequestsimpl

import (
	"context"

	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
)

type FakeService struct {
	ExpectedError         error
	ExpectedUploadRequest uploadrequests.Model
}

func NewFakeService() uploadrequests.Service {
	return &FakeService{}
}

var _ uploadrequests.Service = (*FakeService)(nil)

func (s *FakeService) Create(ctx context.Context, cmd uploadrequests.CreateUploadRequestCommand) (uploadrequests.Model, error) {
	return s.ExpectedUploadRequest, s.ExpectedError
}

func (s *FakeService) GetByID(ctx context.Context, query uploadrequests.GetByIDQuery) (uploadrequests.Model, error) {
	return s.ExpectedUploadRequest, s.ExpectedError
}

func (s *FakeService) UpdateStatus(ctx context.Context, cmd uploadrequests.UpdateStatusCommand) error {
	return s.ExpectedError
}
