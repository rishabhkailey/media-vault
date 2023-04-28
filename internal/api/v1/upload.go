package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/services/media"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type initRequestBody struct {
	FileName  string `json:"fileName"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType,omitempty"`
	Date      int64  `json:"date,omitempty"`
}

func (server *Server) InitChunkUpload(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[InitChunkUpload]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var requestBody initRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logrus.Infof("[InitChunkUpload] invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if len(requestBody.FileName) == 0 || requestBody.Size == 0 {
		logrus.Infof("[InitChunkUpload] invalid request")
		c.Status(http.StatusBadRequest)
		return
	}
	if len(requestBody.MediaType) == 0 {
		requestBody.MediaType = string(mediametadata.TYPE_UNKNOWN)
	}
	if requestBody.Date == 0 {
		requestBody.Date = time.Now().Unix()
	}
	uploadRequest, err := server.UploadRequests.Create(c.Request.Context(), uploadrequests.CreateUploadRequestCommand{
		UserID: userID,
	})
	if err != nil {
		logrus.Errorf("[InitChunkUpload] uploadRequest creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	mediaMetadata, err := server.MediaMetadata.Create(c.Request.Context(), mediametadata.CreateCommand{
		Metadata: mediametadata.Metadata{
			Name: requestBody.FileName,
			Date: time.UnixMilli(requestBody.Date),
			Size: uint64(requestBody.Size),
			Type: requestBody.MediaType,
		},
	})
	if err != nil {
		logrus.Errorf("[InitUpload] media metadata creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	media, err := server.Media.Create(c.Request.Context(), media.CreateMediaCommand{
		UploadRequestID: uploadRequest.ID,
		MetadataID:      mediaMetadata.ID,
	})
	if err != nil {
		logrus.Errorf("[InitUpload] media creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = server.UserMediaBindings.Create(c.Request.Context(), usermediabindings.CreateCommand{
		UserID:  userID,
		MediaID: media.ID,
	})
	if err != nil {
		logrus.Errorf("[InitUpload] UserMediaBindings creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = server.MediaStorage.InitChunkUpload(c.Request.Context(), mediastorage.InitChunkUploadCmd{
		UserID:    userID,
		RequestID: uploadRequest.ID,
		FileName:  media.FileName,
		FileSize:  requestBody.Size,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.InitUpload"}).Errorf("saveUserUploadRequest failed")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requestID": uploadRequest.ID,
	})
}

// type uploadChunkRequestBody struct {
// 	RequestID string `json:"requestID"`
// 	ChunkSize int64  `json:"chunkSize"`
// 	// type?
// 	ChunkData io.Reader `json:"chunkData"`
// 	Index     int64     `json:"index"`
// }

// todo check for multipart request
// todo request size limit
func (server *Server) UploadChunk(c *gin.Context) {
	if c.Request.MultipartForm == nil {
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			logrus.Error("[Server.uploadChunk] parse multpart form failed: %w", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	requestID := c.Request.PostFormValue("requestID")
	if len(requestID) == 0 {
		logrus.Error("[Server.uploadChunk] bad request: requestID param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[InitChunkUpload]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var err error
	var index int64
	{
		value := c.Request.PostFormValue("index")
		if len(value) == 0 {
			logrus.Error("[Server.uploadChunk] bad request: index param missing")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		index, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Error("[Server.uploadChunk] bad request: invalid index %v", value)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	var chunkSize int64
	{
		value := c.Request.PostFormValue("chunkSize")
		if len(value) == 0 {
			logrus.Error("[Server.uploadChunk] bad request: chunkSize param missing")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		chunkSize, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Error("[Server.uploadChunk] bad request: invalid chunkSize %v", value)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	chunkData, _, err := c.Request.FormFile("chunkData")
	if err != nil {
		logrus.Error("[Server.uploadChunk] bad request: chunkData param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	n, err := server.MediaStorage.UploadChunk(c.Request.Context(), mediastorage.UploadChunkCmd{
		UserID:          userID,
		UploadRequestID: requestID,
		Index:           index,
		ChunkSize:       chunkSize,
		Chunk:           chunkData,
	})
	if err != nil {
		logrus.Errorf("[Server.uploadChunk] io.CopyN failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requestID": requestID,
		"uploaded":  n,
	})
}

// thumbnail is required to be of jpeg type only
func (server *Server) UploadThumbnail(c *gin.Context) {
	if c.Request.MultipartForm == nil {
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			logrus.Error("[Server.uploadChunk] parse multpart form failed: %w", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	requestID := c.Request.PostFormValue("requestID")
	if len(requestID) == 0 {
		logrus.Error("[Server.UploadThumbnail] bad request: requestID param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var err error
	var size int64
	{
		value := c.Request.PostFormValue("size")
		if len(value) == 0 {
			logrus.Error("[Server.UploadThumbnail] bad request: size param missing")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		size, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Error("[Server.UploadThumbnail] bad request: invalid size %v", value)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	// todo change the param to thumbnailData? consistency
	thumbnail, _, err := c.Request.FormFile("thumbnail")
	if err != nil {
		logrus.Error("[Server.UploadThumbnail] bad request: chunkData param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[UploadThumbnail]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	media, err := server.Media.GetByUploadRequestID(c.Request.Context(), media.GetByUploadRequestQuery{
		UploadRequestID: requestID,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		logrus.Errorf("[UploadThumbnail]: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = server.MediaStorage.ThumbnailUpload(c.Request.Context(), mediastorage.UploadThumbnailCmd{
		RequestID:  requestID,
		UserID:     userID,
		FileName:   media.FileName,
		FileSize:   size,
		FileReader: thumbnail,
	})
	if err != nil {
		logrus.Errorf("[server.UploadThumbnail] failed to upload file to minio: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = server.MediaMetadata.UpdateThumbnail(c.Request.Context(), mediametadata.UpdateThumbnailCommand{
		Thumbnail: true,
		ID:        media.MetadataID,
	})
	if err != nil {
		logrus.Errorf("[server.UploadThumbnail] failed to update media metadata: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

type finishUploadRequestBody struct {
	RequestID string `json:"requestID"`
	Checksum  string `json:"checksum"`
}

// this is for client to confirm if the upload has finished or not
func (server *Server) FinishChunkUpload(c *gin.Context) {
	var requestBody finishUploadRequestBody
	c.Request.ParseForm()
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logrus.Errorf("[server.FinishUpload] bad request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[FinishChunkUpload]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = server.MediaStorage.FinishChunkUpload(c.Request.Context(), mediastorage.FinishChunkUpload{
		UserID:    userID,
		RequestID: requestBody.RequestID,
		CheckSum:  requestBody.Checksum,
	})
	if err != nil {
		logrus.Errorf("[server.FinishUpload] upload failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	media, err := server.Media.GetMediaWithMetadataByUploadRequestID(c.Request.Context(), media.GetByUploadRequestQuery{
		UploadRequestID: requestBody.RequestID,
	})
	if err != nil {
		logrus.Errorf("[UploadThumbnail]: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	mediaApiData, err := NewMediaApiData(media)
	if err != nil {
		logrus.Errorf("[UploadThumbnail] NewMediaApiData failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mediaApiData)
}
