package authservice

import (
	"context"
	"errors"
	"net/http"
)

type SessionStoreQuery struct {
	Ctx            context.Context
	ResponseWriter http.ResponseWriter
	Request        http.Request
}

type TerminateSessionCmd struct {
	SessionStoreQuery
}

type ValidateUserAccessQuery struct {
	SessionStoreQuery
}

type GetSessionExpireTimeQuery struct {
	SessionStoreQuery
}

type ValidateUserMediaAccessQuery struct {
	SessionStoreQuery
	FileName string
}

var ErrUnauthorized = errors.New("unauthorized")

// todo change scope to unique-id/user and unique-id/admin
// const UserScope = "media-service/user"
// const AdminScope = "media-service/admin"
const UserScope = "user"
const AdminScope = "admin"
