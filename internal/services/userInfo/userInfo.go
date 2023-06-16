package userinfo

import (
	"context"

	userinfo "github.com/rishabhkailey/media-service/internal/store/userInfo"
)

type Service interface {
	GetUserInfo(context.Context, GetUserInfoQuery) (userinfo.UserInfo, error)
	CreateUserInfo(context.Context, CreateUserInfoCmd) (userinfo.UserInfo, error)
}
