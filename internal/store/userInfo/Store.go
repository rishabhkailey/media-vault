package userinfo

import "context"

type Store interface {
	GetByID(ctx context.Context, userID string) (UserInfo, error)
	Insert(ctx context.Context, userID string, encryptionKeyChecksum string, preferedTimezone string) (UserInfo, error)
}
