package v1

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/services/media"
	mediametadata "github.com/rishabhkailey/media-service/internal/services/mediaMetadata"
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
		mediaData, err := NewMediaApiData(media)
		if err != nil {
			logrus.Errorf("[MediaList] NewMediaApiData failed: %w", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		response = append(response, mediaData)
	}
	c.JSON(http.StatusOK, &response)
}

func NewMediaApiData(media media.Model) (MediaApiData, error) {
	var mediaData MediaApiData
	mediaUrl, err := parseMediaURL(media.FileName, false)
	if err != nil {
		return mediaData, err
	}
	mediaData.MediaUrl = mediaUrl
	mediaData.Metadata = media.Metadata.Metadata
	if media.Metadata.Thumbnail {
		thumbnailUrl, err := parseMediaURL(media.FileName, true)
		if err != nil {
			return mediaData, err
		}
		mediaData.ThumbnailUrl = thumbnailUrl
	}
	return mediaData, nil
}

func parseMediaURL(fileName string, thumbnail bool) (string, error) {
	path := "/v1/media"
	if thumbnail {
		path = "/v1/thumbnail"
	}
	return url.JoinPath(path, fileName)
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
	mediaType, err := server.Media.GetTypeByFileName(c.Request.Context(), media.GetTypeByFileNameQuery{
		FileName: fileName,
	})
	if err != nil {
		logrus.Errorf("[GetMedia] get media type failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		logrus.Errorf("[getMedia] failed to write media data to response: %w. expected bytes=%d, written bytes=%d,", err)
		c.Status(http.StatusInternalServerError)
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
	_, err := server.MediaStorage.HttpGetRangeHandler(c.Request.Context(), mediastorage.WriteRangeByFileNameQuery{
		FileName:       fileName,
		Range:          r,
		ResponseWriter: c.Writer,
	})
	if err != nil {
		logrus.Errorf("[getMedia] failed to write media data to response: %w. expected bytes=%d, written bytes=%d,", err)
		c.AbortWithStatus(http.StatusInternalServerError) // this will fail if server has already wrote some data
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
		logrus.Errorf("[getMedia] failed to write thumbnail data to response: %w. expected bytes=%d, written bytes=%d,", err)
		c.AbortWithStatus(http.StatusInternalServerError) // this will fail if server has already wrote some data
		return
	}
}
