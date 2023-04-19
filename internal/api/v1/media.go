package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/services/media"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
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

type MediaListRequestParams struct {
	Page    uint   `form:"page" json:"page,omitempty"`
	PerPage uint   `form:"perPage" json:"perPage,omitempty"`
	OrderBy string `form:"order" json:"order,omitempty"`
	Sort    string `form:"sort" json:"sort,omitempty"`
	// MediaType []string `json:"mediaType,omitempty"`
}

type MediaApiData struct {
	MediaUrl     string `json:"url"`
	ThumbnailUrl string `json:"thumbnail_url"`
	mediametadata.Metadata
}

// todo- ignore upload status failed media
func (server *Server) MediaList(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[GetMediaList]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var requestBody MediaListRequestParams
	if err := c.Bind(&requestBody); err != nil {
		logrus.Infof("[InitUpload] invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	requestBody.initDefaultValues()
	if err := requestBody.validate(); err != nil {
		logrus.Errorf("[MediaList] bad request: %w", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	mediaList, err := server.Media.GetByUserID(c.Request.Context(), media.GetByUserIDQuery{
		UserID:  userID,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
		Offset:  int((requestBody.Page - 1) * requestBody.PerPage),
		Limit:   int(requestBody.PerPage),
	})
	if err != nil {
		logrus.Errorf("[MediaList] db.GetUserMediaList failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response := []MediaApiData{}
	for _, media := range mediaList {
		mediaUrl, err := server.parseMediaURL(media.FileName, false)
		if err != nil {
			logrus.Errorf("[MediaList] error parsing media url: %w", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		mediaData := MediaApiData{
			MediaUrl: mediaUrl,
			Metadata: media.Metadata.Metadata,
		}
		if media.Metadata.Thumbnail {
			thumbnailUrl, err := server.parseMediaURL(media.FileName, true)
			if err != nil {
				logrus.Errorf("[MediaList] error parsing media url: %w", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			mediaData.ThumbnailUrl = thumbnailUrl
		}
		response = append(response, mediaData)
	}
	c.JSON(http.StatusOK, &response)
}

func (server *Server) parseMediaURL(fileName string, thumbnail bool) (string, error) {
	path := "/v1/media"
	if thumbnail {
		path = "/v1/thumbnail"
	}
	return url.JoinPath(path, fileName)
}

func genThumbnailName(fileName string) string {
	return fmt.Sprintf(".thumb-%s", fileName)
}

func (requestBody *MediaListRequestParams) initDefaultValues() {
	if requestBody.Page == 0 {
		requestBody.Page = 1
	}
	if requestBody.PerPage == 0 {
		requestBody.PerPage = MEDIA_API_DEFAULT_PER_PAGE
	}
	if requestBody.PerPage > MEDIA_API_MAX_PER_PAGE {
		requestBody.PerPage = MEDIA_API_MAX_PER_PAGE
	}
	if len(requestBody.OrderBy) == 0 {
		requestBody.OrderBy = usermediabindings.ORDER_BY_MEDIA_CREATION_TIME
	}
	if len(requestBody.Sort) == 0 {
		requestBody.Sort = usermediabindings.SORT_DESCENDING
	}
}

func (requestBody *MediaListRequestParams) validate() error {
	if !utils.Contains(usermediabindings.SUPPORTED_ORDER_BY, requestBody.OrderBy) {
		return fmt.Errorf("invalid orderBy value")
	}
	if requestBody.Sort != usermediabindings.SORT_ASCENDING && requestBody.Sort != usermediabindings.SORT_DESCENDING {
		return fmt.Errorf("invalid sort value")
	}
	return nil
}

func (server *Server) GetMedia(c *gin.Context) {
	fileName := c.Param("fileName")
	if len(fileName) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	rangeHeader := c.Request.Header["Range"]
	var parsedRangeHeader *utils.RangeHeader
	if len(rangeHeader) != 0 && len(rangeHeader[0]) != 0 {
		var err error
		parsedRangeHeader, err = utils.ParseRangeHeader(rangeHeader[0])
		if err != nil {
			logrus.Errorf("[GetMedia] parse range header failed: %w", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	object, err := server.MinioObjectCache.Get(c.Request.Context(), "test", fileName)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	mediaType, err := server.getMediaType(c.Request.Context(), fileName)
	if err != nil {
		logrus.Errorf("[GetMedia] get media type failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// todo support for multiple ranges
	if parsedRangeHeader == nil || len(parsedRangeHeader.Ranges) != 1 {
		server.getMedia(c, object, mediaType)
		return
	}
	server.GetMediaRange(c, parsedRangeHeader.Ranges[0], object, mediaType)
}

func (server *Server) getMedia(c *gin.Context, object *minio.Object, contentType string) {
	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
	}
	_, err = object.Seek(0, 0)
	if err != nil {
		logrus.Errorf("[getMedia] object seek failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// this is for giving client the hint that response is a video file
	contentLength := objInfo.Size
	n, err := io.CopyN(c.Writer, object, objInfo.Size)
	if err != nil || n != objInfo.Size {
		logrus.Errorf("[getMedia] failed to write thumbnail data respose: %w. expected bytes=%d, written bytes=%d,", err, objInfo.Size, n)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", contentType)
	c.Header("Connection", "keep-alive")
	c.Header("Accept-Ranges", "bytes")
	c.Status(http.StatusOK)
}

// todo browsers which don't support range requests
// todo what to do on first request without range
// https://vjs.zencdn.net/v/oceans.mp4 this return a 200 response with content length only?
// if range end not provided
const defaultRangeSize int64 = 1000000 // 1mb
func (server *Server) GetMediaRange(c *gin.Context, r utils.Range, object *minio.Object, contentType string) {

	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
	}
	if objInfo.Size == 0 {
		// todo gracefully handle to this for empty files
		c.Status(http.StatusInternalServerError)
		return
	}
	if r.End == -1 {
		r.End = r.Start + defaultRangeSize
	}
	if r.End > objInfo.Size-1 {
		r.End = objInfo.Size - 1
	}
	contentLength := r.End - r.Start + 1
	// c.SSEvent()
	// todo use of stream?
	logrus.WithField("range", r).Info("request received")
	_, err = object.Seek(r.Start, 0)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// can not change status after we have started writting the response
	c.Status(http.StatusPartialContent)
	// todo request can get stuck in js if it doesn't receive the data equals to content length
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", contentType)
	c.Header("Connection", "keep-alive")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", r.Start, r.End, objInfo.Size))
	c.Header("Accept-Ranges", "bytes")
	n, err := io.CopyN(c.Writer, object, contentLength)
	logrus.WithField("bytes", n).Info("sent")
	if err != nil {
		// todo this will not help i guess, status code set earlier will be sent when we start copying the data
		c.Status(http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
}

func (server *Server) getMediaType(ctx context.Context, fileName string) (mediaType string, err error) {
	mediaType, err = server.RedisStore.GetMediaType(ctx, fileName)
	if err == nil && len(mediaType) != 0 {
		return mediaType, err
	}
	mediaType, err = server.Media.GetTypeByFileName(ctx, media.GetTypeByFileNameQuery{
		FileName: fileName,
	})
	if err != nil {
		if err := server.RedisStore.SetMediaType(ctx, fileName, mediaType); err != nil {
			logrus.Warnf("[getMediaType] redis set media type failed: %w", err)
		}
	}
	return mediaType, err
}

func (server *Server) GetThumbnail(c *gin.Context) {
	fileName := c.Param("fileName")
	if len(fileName) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	thumbnailFile := genThumbnailName(fileName)
	object, err := server.MinioObjectCache.Get(c.Request.Context(), "test", thumbnailFile)
	if err != nil {
		logrus.Errorf("[GetThumbnail] failed to get object: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	objInfo, err := object.Stat()
	if err != nil {
		logrus.Errorf("[GetThumbnail] failed to get object info: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = object.Seek(0, 0)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", 0, objInfo.Size, objInfo.Size))
	n, err := io.CopyN(c.Writer, object, objInfo.Size)
	if err != nil || n != objInfo.Size {
		logrus.Errorf("[GetThumbnail] failed to write thumbnail data respose: %w. expected bytes=%d, written bytes=%d,", err, objInfo.Size, n)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// setting content length and not returning that conent will cause the js fetch to stuck so only set the content length after the data is written
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))
	c.Header("Content-Type", mediametadata.TYPE_IMAGE_JPEG)
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
}
