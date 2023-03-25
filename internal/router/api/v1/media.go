package v1

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	dbservices "github.com/rishabhkailey/media-service/internal/db/services"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

// todo better name
const MAX_MEDIA_PER_PAGE uint = 100
const DEFAULT_MEDIA_PER_PAGE uint = 30

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
	dbservices.Metadata
}

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
	mediaList, err := server.Media.GetUserMediaList(c.Request.Context(), userID, requestBody.OrderBy, requestBody.Sort, int((requestBody.Page-1)*requestBody.PerPage), int(requestBody.PerPage))
	if err != nil {
		logrus.Errorf("[MediaList] db.GetUserMediaList failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response := []MediaApiData{}
	for _, media := range mediaList {
		mediaData := MediaApiData{
			MediaUrl: server.parseMediaURL(media.FileName, false),
			Metadata: media.Metadata.Metadata,
		}
		if media.Metadata.Thumbnail {
			mediaData.ThumbnailUrl = server.parseMediaURL(media.FileName, true)
		}
		response = append(response, mediaData)
	}
	c.JSON(http.StatusOK, &response)
}

func (server *Server) parseMediaURL(fileName string, thumbnail bool) string {
	path := "/v1/media"
	if thumbnail {
		path = "/v1/thumbnail"
	}
	url := url.URL{
		Path:     path,
		RawQuery: fmt.Sprintf("file=%s", fileName),
	}
	return url.String()
}

func genThumbnailName(fileName string) string {
	return fmt.Sprintf(".thumb-%s", fileName)
}

func (requestBody *MediaListRequestParams) initDefaultValues() {
	if requestBody.Page == 0 {
		requestBody.Page = 1
	}
	if requestBody.PerPage == 0 {
		requestBody.PerPage = DEFAULT_MEDIA_PER_PAGE
	}
	if requestBody.PerPage > MAX_MEDIA_PER_PAGE {
		requestBody.PerPage = MAX_MEDIA_PER_PAGE
	}
	if len(requestBody.OrderBy) == 0 {
		requestBody.OrderBy = dbservices.ORDER_BY_UPLOAD_TIME
	}
	if len(requestBody.Sort) == 0 {
		requestBody.Sort = dbservices.SORT_DESCENDING
	}
}

func (requestBody *MediaListRequestParams) validate() error {
	if !utils.Contains(dbservices.SUPPORTED_ORDER_BY, requestBody.OrderBy) {
		return fmt.Errorf("invalid orderBy value")
	}
	if requestBody.Sort != dbservices.SORT_ASCENDING && requestBody.Sort != dbservices.SORT_DESCENDING {
		return fmt.Errorf("invalid sort value")
	}
	return nil
}

func (server *Server) GetMedia(c *gin.Context) {

}
