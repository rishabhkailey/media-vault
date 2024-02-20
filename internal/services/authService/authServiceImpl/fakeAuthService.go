package authserviceimpl

import (
	authservice "github.com/rishabhkailey/media-vault/internal/services/authService"
)

type FakeService struct {
	ExpectedError             error
	ExpectedUserID            string
	ExpectedSessionExpireTime int64
}

func NewFakeService() authservice.Service {
	return &FakeService{}
}

var _ authservice.Service = (*FakeService)(nil)

func (s *FakeService) TerminateSession(cmd authservice.TerminateSessionCmd) error {
	return s.ExpectedError
}

func (s *FakeService) GetSessionExpireTime(query authservice.GetSessionExpireTimeQuery) (int64, error) {
	return s.ExpectedSessionExpireTime, s.ExpectedError
}

func (s *FakeService) ValidateUserAccess(query authservice.ValidateUserAccessQuery, requiredScopes []string, requiredRoles []string) (userID string, err error) {
	return s.ExpectedUserID, s.ExpectedError
}

func (s *FakeService) ValidateUserMediaAccess(query authservice.ValidateUserMediaAccessQuery) (err error) {
	return s.ExpectedError
}

func (s *FakeService) RefreshSession(query authservice.RefreshSessionQuery) (int64, error) {
	return s.ExpectedSessionExpireTime, s.ExpectedError
}
