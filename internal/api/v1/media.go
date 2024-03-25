package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-vault/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-vault/internal/errors"
	"github.com/rishabhkailey/media-vault/internal/services/media"
	mediametadata "github.com/rishabhkailey/media-vault/internal/services/mediaMetadata"
	mediasearch "github.com/rishabhkailey/media-vault/internal/services/mediaSearch"
	mediastorage "github.com/rishabhkailey/media-vault/internal/services/mediaStorage"
	usermediabindings "github.com/rishabhkailey/media-vault/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-vault/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	MEDIA_API_MAX_PER_PAGE                 uint   = 100
	MEDIA_API_DEFAULT_PER_PAGE             uint   = 30
	MEDIA_API_ORDER_BY_UPLOAD_TIME         string = "created_at"
	MEDIA_API_ORDER_BY_MEDIA_CREATION_TIME string = "date"
)

var SUPPORTED_ORDER_BY = []string{MEDIA_API_ORDER_BY_MEDIA_CREATION_TIME, MEDIA_API_ORDER_BY_UPLOAD_TIME}

func (server *Server) GetMedia(c *gin.Context) {
	userID := c.GetString("user_id")
	mediaID, err := strconv.ParseUint(c.Param("media_id"), 10, 64)
	if err != nil {
		c.Error(internalErrors.NewBadRequestError(
			fmt.Errorf("[GetMedia] error parsing mediaID: %w", err),
			"invalid media id",
		))
		return
	}

	belongsToUser, err := server.UserMediaBindings.CheckMediaBelongsToUser(c.Request.Context(), usermediabindings.CheckMediaBelongsToUserQuery{
		MediaID: uint(mediaID),
		UserID:  userID,
	})
	if err != nil {
		c.Error((internalErrors.NewInternalServerError(err)))
		return
	}
	if !belongsToUser {
		c.Error((internalErrors.NewForbiddenError(err)))
		return
	}

	media, err := server.Media.GetByMediaID(c.Request.Context(), media.GetByMediaIDQuery{
		MediaID: uint(mediaID),
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetMedia] db.GetByMediaID failed: %w", err),
			),
		)
		return
	}

	response, err := v1models.NewGetMediaResponse(media)
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetMedia] NewGetMediaResponse failed: %w", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &v1models.FinishUploadResponse{
		GetMediaResponse: response,
	})
}

// todo- ignore upload status failed media
func (server *Server) MediaList(c *gin.Context) {
	userID := c.GetString("user_id")
	if len(userID) == 0 {
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
	var err error
	mediaList, err := server.Media.GetByUserID(c.Request.Context(), media.GetByUserIDQuery{
		UserID:      userID,
		OrderBy:     requestBody.OrderBy,
		Sort:        requestBody.Sort,
		LastMediaID: requestBody.LastMediaID,
		LastDate:    requestBody.LastDate,
		PerPage:     requestBody.PerPage,
	})

	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[MediaList] db.GetUserMediaList failed: %w", err),
			),
		)
		return
	}
	response, err := v1models.NewGetMediaListResponse(mediaList)
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[GetMediaList] create response failed: %w", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &response)
}

func (server *Server) GetMediaFile(c *gin.Context) {
	fileName := c.Param("file_name")
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
		server.getMediaFile(c, fileName, mediaType)
		return
	}
	server.GetMediaFileRange(c, parsedRangeHeader.Ranges[0], fileName, mediaType) // todo support for multiple ranges
}

func (server *Server) getMediaFile(c *gin.Context, fileName string, contentType string) {
	// we can not set headers and status once we have started writting the response
	if contentType != mediametadata.TYPE_UNKNOWN {
		c.Header("Content-Type", contentType)
	}
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	err := server.MediaStorage.HttpGetMediaHandler(c.Request.Context(), mediastorage.HttpGetMediaHandlerQuery{
		FileName:       fileName,
		ResponseWriter: c.Writer,
		Request:        c.Request,
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
func (server *Server) GetMediaFileRange(c *gin.Context, r utils.Range, fileName string, contentType string) {
	c.Status(http.StatusPartialContent)
	if contentType != mediametadata.TYPE_UNKNOWN {
		c.Header("Content-Type", contentType)
	}
	c.Header("Connection", "keep-alive")
	c.Header("Accept-Ranges", "bytes")
	err := server.MediaStorage.HttpGetRangeHandler(c.Request.Context(), mediastorage.HttpGetRangeHandlerQuery{
		FileName:       fileName,
		Range:          r,
		ResponseWriter: c.Writer,
		Request:        c.Request,
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

func (server *Server) GetThumbnailFile(c *gin.Context) {
	fileName := c.Param("file_name")
	if len(fileName) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Header("Content-Type", mediametadata.TYPE_IMAGE_JPEG)
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	err := server.MediaStorage.HttpGetThumbnailHandler(c.Request.Context(), mediastorage.HttpGetThumbnailHandlerQuery{
		FileName:       fileName,
		ResponseWriter: c.Writer,
		Request:        c.Request,
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

func (server *Server) DeleteMedia(c *gin.Context) {
	userID := c.GetString("user_id")
	var requestBody v1models.DeleteMediaRequest
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

	belongsToUser, err := server.UserMediaBindings.CheckMultipleMediaBelongsToUser(c.Request.Context(), usermediabindings.CheckMultipleMediaBelongsToUserQuery{
		UserID:   userID,
		MediaIDs: requestBody.MediaIDs,
	})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] error checking user access: %w", err),
		))
		return
	}
	if !belongsToUser {
		c.Error(internalErrors.ErrForbidden)
		return
	}

	mediaToDelete, err := server.Media.GetByMediaIDs(c.Request.Context(),
		media.GetByMediaIDsQuery{
			MediaIDs: requestBody.MediaIDs,
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] select media query failed: %w", err),
		))
		return
	}

	// delete db entries
	_, _, _, _, err = server.Media.CascadeDeleteMany(c.Request.Context(),
		media.DeleteManyCommand{
			MediaIDs: requestBody.MediaIDs,
			UserID:   userID,
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] failed: %w", err),
		))
		return
	}

	// todo use 422 Unprocessable Entity status code?
	// delete from search service
	searchId, err := server.MediaSearch.DeleteMany(c.Request.Context(),
		mediasearch.DeleteManyCommand{
			MediaIDs: requestBody.MediaIDs,
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] delete from search service failed, this will cause discrepancies in search results. document id = %d: %w", searchId, err),
		))
		return
	}

	var deleteMediaFileCmd mediastorage.DeleteManyCommand
	for _, media := range mediaToDelete {
		deleteMediaFileCmd.DeleteCmds = append(deleteMediaFileCmd.DeleteCmds, mediastorage.DeleteOneCommand{
			FileName:     media.FileName,
			HasThumbnail: media.Metadata.Thumbnail,
		})
	}

	// delete media file
	failedFileNames, errors := server.MediaStorage.DeleteMany(c.Request.Context(),
		deleteMediaFileCmd)
	if len(errors) != 0 {
		for index, failedFileNames := range failedFileNames {
			logrus.Warnf("[DeleteMedia] failed to delete %s file: %v", failedFileNames, errors[index])
		}
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] failed to delete %d files from storage", len(errors)),
		))
		return
	}

	c.Status(http.StatusOK)
}

func (server *Server) DeleteSingleMedia(c *gin.Context) {
	userID := c.GetString("user_id")
	mediaID, err := strconv.ParseUint(c.Param("media_id"), 10, 64)
	if err != nil {
		c.Error(internalErrors.NewBadRequestError(
			fmt.Errorf("[DeleteMedia] error parsing mediaID: %w", err),
			"invalid media id",
		))
		return
	}

	belongsToUser, err := server.UserMediaBindings.CheckMediaBelongsToUser(c.Request.Context(), usermediabindings.CheckMediaBelongsToUserQuery{
		UserID:  userID,
		MediaID: uint(mediaID),
	})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] error checking user access: %w", err),
		))
		return
	}
	if !belongsToUser {
		c.Error(internalErrors.ErrForbidden)
		return
	}

	mediaToDelete, err := server.Media.GetByMediaID(c.Request.Context(),
		media.GetByMediaIDQuery{
			MediaID: uint(mediaID),
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] select media query failed: %w", err),
		))
		return
	}

	// delete db entries
	err = server.Media.CascadeDeleteOne(c.Request.Context(),
		media.DeleteOneCommand{
			ID:         uint(mediaID),
			UserID:     userID,
			MetadataID: mediaToDelete.Metadata.ID,
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] failed: %w", err),
		))
		return
	}

	// delete from search service
	searchId, err := server.MediaSearch.DeleteOne(c.Request.Context(),
		mediasearch.DeleteOneCommand{
			MediaID: uint(mediaID),
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] delete from search service failed, this will cause discrepancies in search results. document id = %d: %w", searchId, err),
		))
		return
	}

	// delete media file
	err = server.MediaStorage.DeleteOne(c.Request.Context(),
		mediastorage.DeleteOneCommand{
			FileName:     mediaToDelete.FileName,
			HasThumbnail: mediaToDelete.Metadata.Thumbnail,
		})
	if err != nil {
		c.Error(internalErrors.NewInternalServerError(
			fmt.Errorf("[DeleteMedia] delete from storage failed, filename = %s: %w", mediaToDelete.FileName, err),
		))
		return
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
// 	userID, ok := c.Keys["user_id"].(string)
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
