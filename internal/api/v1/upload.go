package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/services/media/mediaimpl"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"gorm.io/gorm"
)

func (server *Server) InitChunkUpload(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[InitChunkUpload]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.InitChunkUploadRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[InitChunkUpload] invalid request: %v", err),
				"bad request",
			),
		)
		return
	}
	requestBody, err := v1models.ValidateInitChunkUploadRequest(requestBody)
	if err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[InitChunkUpload] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	var uploadRequest uploadrequests.Model
	var uploadingMedia media.Model
	{
		tx := server.Services.CreateTransaction()
		uploadRequest, err = server.UploadRequests.Create(c.Request.Context(), uploadrequests.CreateUploadRequestCommand{
			UserID: userID,
		})
		if err != nil {
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[InitChunkUpload] uploadRequest creation failed: %w", err),
				),
			)
			return
		}
		mediaMetadata, err := server.MediaMetadata.WithTransaction(tx).Create(c.Request.Context(), mediametadata.CreateCommand{
			Metadata: mediametadata.Metadata{
				Name: requestBody.FileName,
				Date: time.UnixMilli(requestBody.Date),
				Size: uint64(requestBody.Size),
				Type: requestBody.MediaType,
			},
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[InitUpload] media metadata creation failed: %w", err),
				),
			)
			return
		}
		uploadingMedia, err = server.Media.WithTransaction(tx).Create(c.Request.Context(), media.CreateMediaCommand{
			UploadRequestID: uploadRequest.ID,
			MetadataID:      mediaMetadata.ID,
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[InitUpload] media creation failed: %w", err),
				),
			)
			return
		}
		_, err = server.UserMediaBindings.WithTransaction(tx).Create(c.Request.Context(), usermediabindings.CreateCommand{
			UserID:  userID,
			MediaID: uploadingMedia.ID,
		})
		if err != nil {
			tx.Rollback()
			c.Error(
				internalErrors.NewInternalServerError(
					fmt.Errorf("[InitUpload] UserMediaBindings creation failed: %w", err),
				),
			)
			return
		}
		tx.Commit()
	}
	err = server.MediaStorage.InitChunkUpload(c.Request.Context(), mediastorage.InitChunkUploadCmd{
		UserID:    userID,
		RequestID: uploadRequest.ID,
		FileName:  uploadingMedia.FileName,
		FileSize:  requestBody.Size,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[InitUpload] InitChunkUpload failed: %w", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &v1models.InitChunkUploadResponse{
		RequestID: uploadRequest.ID,
		FileName:  uploadingMedia.FileName,
	})
}

func (server *Server) UploadChunk(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[InitChunkUpload]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.UploadChunkRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[Server.uploadChunk] bad request: %w", err),
				"bad request",
			),
		)
		return
	}
	chunkFile, err := requestBody.ChunkData.Open()
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Server.uploadChunk] chunk data file open failed: %w", err),
			),
		)
		return
	}

	n, err := server.MediaStorage.UploadChunk(c.Request.Context(), mediastorage.UploadChunkCmd{
		UserID:          userID,
		UploadRequestID: requestBody.RequestID,
		Index:           *requestBody.Index,
		ChunkSize:       requestBody.ChunkSize,
		Chunk:           chunkFile,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Server.uploadChunk] io.CopyN failed: %v", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &v1models.UploadChunkResponse{
		RequestID: requestBody.RequestID,
		Uploaded:  n,
	})
}

// thumbnail is required to be of jpeg type only
func (server *Server) UploadThumbnail(c *gin.Context) {
	var requestBody v1models.UploadThumbnailRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[Server.UploadThumbnail] bad request: %w", err),
				"bad request",
			),
		)
		return
	}
	thumbnailFile, err := requestBody.Thumbnail.Open()
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Server.UploadThumbnail] thumbnail file open failed: %w", err),
			),
		)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[UploadThumbnail]: empty userID"),
			),
		)
		return
	}
	media, err := server.Media.GetByUploadRequestID(c.Request.Context(), media.GetByUploadRequestQuery{
		UploadRequestID: requestBody.RequestID,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Error(internalErrors.ErrForbidden)
		return
	}
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[UploadThumbnail] get media by upload request id failed: %w", err),
			),
		)
		return
	}
	err = server.MediaStorage.ThumbnailUpload(c.Request.Context(), mediastorage.UploadThumbnailCmd{
		RequestID:  requestBody.RequestID,
		UserID:     userID,
		FileName:   media.FileName,
		FileSize:   requestBody.Size,
		FileReader: thumbnailFile,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[server.UploadThumbnail] failed to upload file to minio: %w", err),
			),
		)
		return
	}
	err = server.MediaMetadata.UpdateThumbnail(c.Request.Context(), mediametadata.UpdateThumbnailCommand{
		Thumbnail: true,
		ID:        media.MetadataID,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[server.UploadThumbnail] failed to update media metadata: %w", err),
			),
		)
		return
	}
	c.Status(http.StatusOK)
}

// this is for client to confirm if the upload has finished or not
func (server *Server) FinishChunkUpload(c *gin.Context) {
	var requestBody v1models.FinishUploadRequest
	err := c.Bind(&requestBody)
	if err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[server.FinishUpload] bad request: %w", err),
				"bad request",
			),
		)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[FinishChunkUpload]: empty userID"),
			),
		)
		return
	}
	err = server.MediaStorage.FinishChunkUpload(c.Request.Context(), mediastorage.FinishChunkUpload{
		UserID:    userID,
		RequestID: requestBody.RequestID,
		CheckSum:  requestBody.Checksum,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[server.FinishUpload] upload failed: %v", err),
			),
		)
		return
	}
	uploadedMedia, err := server.Media.GetMediaWithMetadataByUploadRequestID(c.Request.Context(), media.GetByUploadRequestQuery{
		UploadRequestID: requestBody.RequestID,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[FinishChunkUpload]: %w", err),
			),
		)
		return
	}
	index, err := mediasearch.MediaToMeiliSearchMediaIndex([]media.Model{uploadedMedia}, userID)
	if err != nil {
		// todo fail request on search index creation fail?
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[FinishChunkUpload]: %w", err),
			),
		)
		return
	}
	_, err = server.MediaSearch.CreateMany(c.Request.Context(), index)
	if err != nil {
		// todo fail request on search index creation fail?
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[FinishChunkUpload]: %w", err),
			),
		)
		return
	}

	mediaApiData, err := mediaimpl.NewGetMediaQueryResultItem(uploadedMedia)
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[UploadThumbnail] NewMediaApiData failed: %w", err),
			),
		)
		return
	}
	c.JSON(http.StatusOK, &v1models.FinishUploadResponse{
		GetMediaQueryResultItem: mediaApiData,
	})
}
