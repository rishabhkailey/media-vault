package db

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
)

type UserSession struct {
	UserID            string
	UserScope         string
	HasUserScope      bool
	SessionExpireTime int64 // unix
}

func GetUserSession(c *gin.Context) (*UserSession, error) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		return nil, fmt.Errorf("[getUserSession]: session start failed")
	}
	var hasUserScope bool
	if value, ok := store.Get("hasUserScope"); ok {
		hasUserScope, _ = value.(bool)
	}
	var sessionExpireTime int64
	if value, ok := store.Get("sessionExpireTime"); ok {
		// ignoring errors here because on error the default sessionExpirteTime will be 0 which is not valid
		svalue, _ := value.(string)
		if len(svalue) != 0 {
			sessionExpireTime, _ = strconv.ParseInt(svalue, 10, 64)
		}
	}
	var userID string
	if value, ok := store.Get("userID"); ok {
		userID, _ = value.(string)
	}
	var userScope string
	if value, ok := store.Get("userScope"); ok {
		userScope, _ = value.(string)
	}

	return &UserSession{
		UserID:            userID,
		UserScope:         userScope,
		HasUserScope:      hasUserScope,
		SessionExpireTime: sessionExpireTime,
	}, nil
}

func SetUserSession(c *gin.Context, userSession UserSession) error {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		return fmt.Errorf("[getUserSession]: session start failed")
	}
	store.Set("userID", userSession.UserID)
	store.Set("userScope", userSession.UserScope)
	store.Set("hasUserScope", userSession.HasUserScope)
	// we will store the unix time as string, as the values are stored in redis we looses the actual type for the value
	// int64 is converted to float64 when jsonunmarshall is called check these 2 issues for better understanding https://github.com/json-iterator/go/issues/351 and https://github.com/json-iterator/go/issues/145
	store.Set("sessionExpireTime", strconv.FormatInt(userSession.SessionExpireTime, 10))
	return store.Save()
}

func GetUserFileAccessFromSession(c *gin.Context, fileName, userID string) (fileBelongsToUser bool, err error) {
	fileBelongsToUser = false
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		return false, err
	}
	key := fmt.Sprintf("%s:%s", fileName, userID)
	if value, ok := store.Get(key); ok {
		fileBelongsToUser, _ = value.(bool)
	}
	return fileBelongsToUser, nil
}

func SetUserFileAccessInSession(c *gin.Context, fileName, userID string) error {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", fileName, userID)
	store.Set(key, true)
	return store.Save()
}

func DeleteUserSession(c *gin.Context) error {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		return err
	}
	store.Delete("userID")
	store.Delete("userScope")
	store.Delete("hasUserScope")
	store.Delete("sessionExpireTime")
	return store.Save()
}
