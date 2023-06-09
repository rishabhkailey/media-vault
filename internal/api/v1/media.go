package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/rishabhkailey/media-service/internal/services/media"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	MEDIA_API_MAX_PER_PAGE                 uint   = 100
	MEDIA_API_DEFAULT_PER_PAGE             uint   = 30
	MEDIA_API_ORDER_BY_UPLOAD_TIME         string = "created_at"
	MEDIA_API_ORDER_BY_MEDIA_CREATION_TIME string = "date"
)

var SUPPORTED_ORDER_BY = []string{MEDIA_API_ORDER_BY_MEDIA_CREATION_TIME, MEDIA_API_ORDER_BY_UPLOAD_TIME}

// todo- ignore upload status failed media
func (server *Server) MediaList(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[MediaList]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.GetMediaListRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[MediaList] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[MediaList] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	var response v1models.GetMediaListResponse
	var err error
	response, err = server.Media.GetByUserID(c.Request.Context(), media.GetByUserIDQuery{
		UserID:  userID,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[MediaList] db.GetUserMediaList failed: %w", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &response)
}

func (server *Server) GetMedia(c *gin.Context) {
	fileName := c.Param("fileName")
	if len(fileName) == 0 {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("fileName param missing"),
				"missing file name",
			),
		)
		return
	}
	rangeHeader := c.Request.Header["Range"]
	var parsedRangeHeader *utils.RangeHeader
	if len(rangeHeader) != 0 && len(rangeHeader[0]) != 0 {
		var err error
		parsedRangeHeader, err = utils.ParseRangeHeader(rangeHeader[0])
		if err != nil {
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[GetMedia] parse range header failed: %w", err),
				),
			)
			return
		}
	}
	mediaType, err := server.Media.GetTypeByFileName(c.Request.Context(), media.GetTypeByFileNameQuery{
		FileName: fileName,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetMedia] get media type failed: %w", err),
			),
		)
		return
	}
	if parsedRangeHeader == nil || len(parsedRangeHeader.Ranges) == 0 {
		server.getMedia(c, fileName, mediaType)
		return
	}
	server.GetMediaRange(c, parsedRangeHeader.Ranges[0], fileName, mediaType) // todo support for multiple ranges
}

func (server *Server) getMedia(c *gin.Context, fileName string, contentType string) {
	// we can not set headers and status once we have started writting the response
	c.Header("Content-Type", contentType)
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	_, err := server.MediaStorage.HttpGetMediaHandler(c.Request.Context(), mediastorage.HttpGetMediaHandlerQuery{
		FileName:       fileName,
		ResponseWriter: c.Writer,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[getMedia] http handler returned error: %w", err),
			),
		)
		return
	}
}

// todo browsers which don't support range requests
// todo what to do on first request without range
// https://vjs.zencdn.net/v/oceans.mp4 this return a 200 response with content length only?
func (server *Server) GetMediaRange(c *gin.Context, r utils.Range, fileName string, contentType string) {
	c.Status(http.StatusPartialContent)
	c.Header("Content-Type", contentType)
	c.Header("Connection", "keep-alive")
	c.Header("Accept-Ranges", "bytes")
	_, err := server.MediaStorage.HttpGetRangeHandler(c.Request.Context(), mediastorage.HttpGetRangeHandlerQuery{
		FileName:       fileName,
		Range:          r,
		ResponseWriter: c.Writer,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetMediaRange] http handler returned error: %w", err),
			),
		)
		return
	}
}

func (server *Server) GetThumbnail(c *gin.Context) {
	fileName := c.Param("fileName")
	if len(fileName) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Header("Content-Type", mediametadata.TYPE_IMAGE_JPEG)
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	_, err := server.MediaStorage.HttpGetThumbnailHandler(c.Request.Context(), mediastorage.HttpGetThumbnailHandlerQuery{
		FileName:       fileName,
		ResponseWriter: c.Writer,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetThumbnail] http handler returned error: %w", err),
			),
		)
		return
	}
}

// todo a child method returning error?
// multiple delete in same API doesn't make sense
// if 1 delete fails the complete request will also fail and as we are using transaction for multiple services it will increase the complexity
func (server *Server) DeleteMedia(c *gin.Context) {
	mediaIdParam := c.Param("mediaID")
	if len(mediaIdParam) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	mediaID, err := strconv.ParseUint(mediaIdParam, 10, 64)
	if err != nil {
		c.Error(
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[DeleteMedia] error parsing mediaID: %w", err),
				),
			),
		)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[DeleteMedia]: empty userID"),
			),
		)
		return
	}
	belongsToUser, err := server.UserMediaBindings.CheckMediaBelongsToUser(c.Request.Context(), usermediabindings.CheckMediaBelongsToUserQuery{
		UserID:  userID,
		MediaID: uint(mediaID),
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[DeleteMedia] error checking user access: %w", err),
			),
		)
		return
	}
	if !belongsToUser {
		c.Error(internalErrors.ErrForbidden)
		return
	}
	{
		deletingMedia, err := server.Media.GetByMediaID(c.Request.Context(), media.GetByMediaIDQuery{
			MediaID: uint(mediaID),
		})
		if err != nil {
			c.Error(
				c.Error(
					internalErrors.NewInternalServerError(
						fmt.Errorf("[DeleteMedia] error while getting media: %w", err),
					),
				),
			)
			return
		}

		tx := server.Services.CreateTransaction()
		err = server.Media.WithTransaction(tx).DeleteOne(c.Request.Context(), media.DeleteOneCommand{
			ID: deletingMedia.ID,
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				c.Error(
					internalErrors.NewInternalServerError(
						fmt.Errorf("[DeleteMedia] error while deleting media: %w", err),
					),
				),
			)
			return
		}

		err = server.MediaMetadata.WithTransaction(tx).DeleteOne(c.Request.Context(), mediametadata.DeleteOneCommand{
			ID: deletingMedia.Metadata.ID,
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				c.Error(
					internalErrors.NewInternalServerError(
						fmt.Errorf("[DeleteMedia] error while deleting media metadata: %w", err),
					),
				),
			)
			return
		}

		_, err = server.MediaSearch.DeleteOne(c.Request.Context(), mediasearch.DeleteOneCommand{
			MediaID: deletingMedia.ID,
		})
		if err != nil {
			// this should not cause much trouble
			logrus.Warnf("[DeleteMedia] delete search document failed: %v", err)
		}
		err = server.MediaStorage.DeleteOne(c.Request.Context(), mediastorage.DeleteOneCommand{
			FileName:     deletingMedia.FileName,
			HasThumbnail: deletingMedia.Metadata.Thumbnail,
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				c.Error(
					internalErrors.NewInternalServerError(
						fmt.Errorf("[DeleteMedia] error while deleting media from storge: %w", err),
					),
				),
			)
			return
		}
		tx.Commit()
	}
	c.Status(http.StatusOK)
}

// func (server *Server) DeleteMultipleMedia(c *gin.Context) {
// 	mediaIdParam := c.Param("mediaID")
// 	if len(mediaIdParam) == 0 {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	mediaID, err := strconv.ParseUint(mediaIdParam, 10, 64)
// 	if err != nil {
// 		c.Error(
// 			c.Error(
// 				internalErrors.NewInternalServerError(
// 					fmt.Errorf("[DeleteMedia] error parsing mediaID: %w", err),
// 				),
// 			),
// 		)
// 		return
// 	}
// 	userID, ok := c.Keys["userID"].(string)
// 	if !ok || len(userID) == 0 {
// 		c.Error(
// 			internalErrors.NewInternalServerError(
// 				fmt.Errorf("[DeleteMedia]: empty userID"),
// 			),
// 		)
// 		return
// 	}
// 	belongsToUser, err := server.UserMediaBindings.CheckMediaBelongsToUser(c.Request.Context(), usermediabindings.CheckMediaBelongsToUserQuery{
// 		UserID:  userID,
// 		MediaID: uint(mediaID),
// 	})
// 	if err != nil {
// 		c.Error(
// 			c.Error(
// 				internalErrors.NewInternalServerError(
// 					fmt.Errorf("[DeleteMedia] error checking user access: %w", err),
// 				),
// 			),
// 		)
// 		return
// 	}
// 	if !belongsToUser {
// 		c.Error(internalErrors.ErrForbidden)
// 		return
// 	}
// 	{
// 		deletingMedia, err := server.Media.GetByMediaID(c.Request.Context(), media.GetByMediaIDQuery{
// 			MediaID: uint(mediaID),
// 		})
// 		if err != nil {
// 			c.Error(
// 				c.Error(
// 					internalErrors.NewInternalServerError(
// 						fmt.Errorf("[DeleteMedia] error while getting media: %w", err),
// 					),
// 				),
// 			)
// 			return
// 		}

// 		tx := server.Services.CreateTransaction()
// 		err = server.MediaMetadata.WithTransaction(tx).DeleteOne(c.Request.Context(), mediametadata.DeleteOneCommand{
// 			ID: deletingMedia.Metadata.ID,
// 		})
// 		if err != nil {
// 			tx.Rollback()
// 			c.Error(
// 				c.Error(
// 					internalErrors.NewInternalServerError(
// 						fmt.Errorf("[DeleteMedia] error while deleting media metadata: %w", err),
// 					),
// 				),
// 			)
// 			return
// 		}

// 		err = server.Media.WithTransaction(tx).DeleteOne(c.Request.Context(), media.DeleteOneCommand{
// 			ID: deletingMedia.ID,
// 		})
// 		if err != nil {
// 			tx.Rollback()
// 			c.Error(
// 				c.Error(
// 					internalErrors.NewInternalServerError(
// 						fmt.Errorf("[DeleteMedia] error while deleting media: %w", err),
// 					),
// 				),
// 			)
// 			return
// 		}
// 		_, err = server.MediaSearch.DeleteOne(c.Request.Context(), mediasearch.DeleteOneCommand{
// 			MediaID: deletingMedia.ID,
// 		})
// 		if err != nil {
// 			// this should not cause much trouble
// 			logrus.Warnf("[DeleteMedia] delete search document failed: %v", err)
// 		}
// 		err = server.MediaStorage.DeleteOne(c.Request.Context(), mediastorage.DeleteOneCommand{
// 			FileName:     deletingMedia.FileName,
// 			HasThumbnail: deletingMedia.Metadata.Thumbnail,
// 		})
// 		if err != nil {
// 			tx.Rollback()
// 			c.Error(
// 				c.Error(
// 					internalErrors.NewInternalServerError(
// 						fmt.Errorf("[DeleteMedia] error while deleting media from storge: %w", err),
// 					),
// 				),
// 			)
// 			return
// 		}
// 		tx.Commit()
// 	}
// 	c.Status(http.StatusOK)
// }
