package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	userinfo "github.com/rishabhkailey/media-service/internal/services/userInfo"
)

func (s *Server) GetUserInfo(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}

	userInfo, err := s.Services.UserInfoService.GetUserInfo(c.Request.Context(), userinfo.GetUserInfoQuery{
		UserID: userID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.PostUserInfoResponse{
		ID:                    userInfo.ID,
		PreferedTimeZone:      userInfo.PreferedTimeZone,
		StorageUsage:          userInfo.StorageUsage,
		EncryptionKeyChecksum: userInfo.EncryptionKeyChecksum,
	})
}

func (s *Server) PostUserInfo(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.PostUserInfoRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[CreateAlbum] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[CreateAlbum] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	userInfo, err := s.Services.UserInfoService.CreateUserInfo(c.Request.Context(), userinfo.CreateUserInfoCmd{
		UserID:                userID,
		PreferedTimeZone:      requestBody.PreferedTimeZone,
		EncryptionKeyChecksum: requestBody.EncryptionKeyChecksum,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.PostUserInfoResponse{
		ID:                    userInfo.ID,
		PreferedTimeZone:      userInfo.PreferedTimeZone,
		StorageUsage:          userInfo.StorageUsage,
		EncryptionKeyChecksum: userInfo.EncryptionKeyChecksum,
	})
}
