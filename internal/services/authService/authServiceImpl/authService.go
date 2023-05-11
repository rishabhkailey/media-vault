package authserviceimpl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-session/session/v3"
	"github.com/rishabhkailey/media-service/internal/auth"
	authservice "github.com/rishabhkailey/media-service/internal/services/authService"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/sirupsen/logrus"
)

type Service struct {
	maxSessionExpireTime time.Duration
	oidcClient           auth.OidcClient
	userMediaBindings    usermediabindings.Service
}

func NewService(oidcClient auth.OidcClient, userMediaBindings usermediabindings.Service, maxSessionExpireTime time.Duration) (authservice.Service, error) {
	return &Service{
		maxSessionExpireTime: maxSessionExpireTime,
		oidcClient:           oidcClient,
		userMediaBindings:    userMediaBindings,
	}, nil
}

var _ authservice.Service = (*Service)(nil)

func (s *Service) TerminateSession(cmd authservice.TerminateSessionCmd) error {
	return session.Destroy(cmd.Ctx, cmd.ResponseWriter, &cmd.Request)
}

func (s *Service) GetSessionExpireTime(query authservice.GetSessionExpireTimeQuery) (sessionExpireTime int64, err error) {
	store, err := session.Start(query.Ctx, query.ResponseWriter, &query.Request)
	if err != nil {
		err = fmt.Errorf("[AuthService.GetSessionExpireTime] session start failed: %w", err)
		return
	}
	if value, ok := store.Get("sessionExpireTime"); ok {
		// ignoring errors here because on error the default sessionExpirteTime will be 0 which is not valid
		svalue, _ := value.(string)
		if len(svalue) != 0 {
			sessionExpireTime, _ = strconv.ParseInt(svalue, 10, 64)
		}
	}
	if sessionExpireTime < time.Now().Unix() {
		err = fmt.Errorf("[AuthService.GetSessionExpireTime] session not found")
		return
	}
	return
}

func (s *Service) ValidateUserAccess(query authservice.ValidateUserAccessQuery, scopes []string) (userID string, err error) {
	userID, userScope, err := s.getUserScope(query.Ctx, query.Request, query.ResponseWriter)
	if err != nil {
		return
	}
	if len(userID) == 0 {
		err = fmt.Errorf("[AuthService.ValidateUserAccess]: empty userID")
		return
	}
	if !auth.ValidateScope(userScope, scopes) {
		err = authservice.ErrUnauthorized
	}
	return
}

func (s *Service) getUserScope(ctx context.Context, r http.Request, w http.ResponseWriter) (userID, userScope string, err error) {
	userID, userScope, err = s.getUserScopeFromSession(ctx, r, w)
	if err == nil {
		return
	}
	logrus.Debug("[AuthService.getUserScope] unable to get user scope from session: %w", err)
	userID, userScope, expireTime, err := s.getUserScopeFromOidcProvider(ctx, r, w)
	if err == nil {
		err := s.saveUserScopeInSession(ctx, r, w, userID, userScope, expireTime)
		if err != nil {
			logrus.Warnf("[AuthService.getUserScope] unable to save session: %w", err)
		}
	}
	return
}

func (s *Service) getUserScopeFromSession(ctx context.Context, r http.Request, w http.ResponseWriter) (userID, userScope string, err error) {
	store, err := session.Start(ctx, w, &r)
	if err != nil {
		err = fmt.Errorf("[AuthService.getUserScopeFromSession] session start failed: %w", err)
		return
	}
	var sessionExpireTime int64
	if value, ok := store.Get("sessionExpireTime"); ok {
		// ignoring errors here because on error the default sessionExpirteTime will be 0 which is not valid
		svalue, _ := value.(string)
		if len(svalue) != 0 {
			sessionExpireTime, _ = strconv.ParseInt(svalue, 10, 64)
		}
	}
	if sessionExpireTime < time.Now().Unix() {
		err = fmt.Errorf("[AuthService.getUserScopeFromSession] session not found")
		return
	}
	if value, ok := store.Get("userID"); ok {
		userID, _ = value.(string)
	}
	if value, ok := store.Get("userScope"); ok {
		userScope, _ = value.(string)
	}
	return
}

func (s *Service) getUserScopeFromOidcProvider(ctx context.Context, r http.Request, w http.ResponseWriter) (userID, userScope string, tokenExpireTime int64, err error) {
	token, ok := auth.GetBearerToken(&r)
	if !ok {
		err = authservice.ErrUnauthorized
		return
	}
	tokenInfo, err := s.oidcClient.IntrospectToken(token)
	if err != nil {
		err = fmt.Errorf("[AuthService.getUserScopeFromOidcProvider] token interospection failed: %w", err)
		return
	}
	userScope = tokenInfo.Scope
	userID = tokenInfo.Subject
	tokenExpireTime = tokenInfo.ExpireTime
	return
}

func (s *Service) saveUserScopeInSession(ctx context.Context, r http.Request, w http.ResponseWriter, userID, userScope string, tokenExpireTime int64) (err error) {
	store, err := session.Start(ctx, w, &r)
	if err != nil {
		return fmt.Errorf("[AuthService.saveUserScopeInSession] session start failed: %w", err)
	}
	store.Set("userID", userID)
	store.Set("userScope", userScope)
	expireTime := time.Now().Add(s.maxSessionExpireTime)
	if expireTime.Unix() > tokenExpireTime {
		expireTime = time.Unix(tokenExpireTime, 0)
	}
	// we will store the unix time as string, as the values are stored in redis we looses the actual type for the value
	// int64 is converted to float64 when jsonunmarshall is called, check these 2 issues for better understanding https://github.com/json-iterator/go/issues/351 and https://github.com/json-iterator/go/issues/145
	store.Set("sessionExpireTime", strconv.FormatInt(expireTime.Unix(), 10))
	err = store.Save()
	if err != nil {
		return fmt.Errorf("[AuthService.saveUserScopeInSession] session save failed: %w", err)
	}
	return
}

func (s *Service) ValidateUserMediaAccess(query authservice.ValidateUserMediaAccessQuery) (err error) {
	userID, _, err := s.getUserScope(query.Ctx, query.Request, query.ResponseWriter)
	if err != nil {
		return
	}
	if len(userID) == 0 {
		err = fmt.Errorf("[AuthService.ValidateUserMediaAccessFromSession]: empty userID")
		return
	}
	err = s.validateUserMediaAccessFromSession(query.Ctx, query.Request, query.ResponseWriter, query.FileName, userID)
	if err == nil || errors.Is(err, authservice.ErrUnauthorized) {
		return
	}
	fileBelongsTofileBelongsToUser, err := s.userMediaBindings.CheckFileBelongsToUser(query.Ctx, usermediabindings.CheckFileBelongsToUserQuery{
		UserID:   userID,
		FileName: query.FileName,
	})
	if err != nil {
		return
	}
	if err := s.SaveUserMediaAccessInSession(query.Ctx, query.Request, query.ResponseWriter, query.FileName, userID, fileBelongsTofileBelongsToUser); err != nil {
		logrus.Warnf("[AuthService.ValidateUserMediaAccess]: saving user access to session failed")
	}
	if !fileBelongsTofileBelongsToUser {
		err = authservice.ErrUnauthorized
	}
	return
}

func userMediaAccessSessionKey(userID, fileName string) string {
	return fmt.Sprintf("%s:%s", fileName, userID)
}

func (s *Service) validateUserMediaAccessFromSession(ctx context.Context, r http.Request, w http.ResponseWriter, fileName string, userID string) (err error) {
	store, err := session.Start(ctx, w, &r)
	if err != nil {
		return err
	}
	fileBelongsToUser := false
	key := userMediaAccessSessionKey(fileName, userID)
	var ok bool
	var value interface{}
	if value, ok = store.Get(key); ok {
		fileBelongsToUser, ok = value.(bool)
	}
	if !ok {
		return fmt.Errorf("[AuthService.validateUserMediaAccessFromSession]: not found")
	}
	if !fileBelongsToUser {
		err = authservice.ErrUnauthorized
	}
	return
}

func (s *Service) SaveUserMediaAccessInSession(ctx context.Context, r http.Request, w http.ResponseWriter, fileName string, userID string, value bool) (err error) {
	store, err := session.Start(ctx, w, &r)
	if err != nil {
		return fmt.Errorf("[AuthService.saveUserScopeInSession] session start failed: %w", err)
	}
	key := userMediaAccessSessionKey(fileName, userID)
	store.Set(key, value)
	return store.Save()
}
