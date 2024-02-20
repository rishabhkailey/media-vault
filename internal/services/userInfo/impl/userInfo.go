package userinfoimpl

import (
	"context"
	"errors"
	"net/http"

	internalErrors "github.com/rishabhkailey/media-vault/internal/errors"
	userinfo "github.com/rishabhkailey/media-vault/internal/services/userInfo"
	"github.com/rishabhkailey/media-vault/internal/store"
	userinfoStore "github.com/rishabhkailey/media-vault/internal/store/userInfo"
	"gorm.io/gorm"
)

type Service struct {
	store store.Store
}

var _ userinfo.Service = (*Service)(nil)

func NewService(store store.Store) (userinfo.Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) GetUserInfo(ctx context.Context, query userinfo.GetUserInfoQuery) (userinfoStore.UserInfo, error) {
	userinfo, err := s.store.UserInfoStore.GetByID(ctx, query.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = internalErrors.NewNotFoundError(err)
	}
	return userinfo, err
}

func (s *Service) CreateUserInfo(ctx context.Context, cmd userinfo.CreateUserInfoCmd) (userInfo userinfoStore.UserInfo, err error) {
	_, err = s.store.UserInfoStore.GetByID(ctx, cmd.UserID)
	if err == nil {
		return userInfo, internalErrors.New(http.StatusConflict, "user info already exists", "already exists")
	}
	return s.store.UserInfoStore.Insert(ctx, cmd.UserID, cmd.EncryptionKeyChecksum, cmd.PreferedTimeZone)
}
