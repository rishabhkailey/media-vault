package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type SearchRequestParams struct {
	Query   string `form:"query" json:"query"`
	Page    uint   `form:"page" json:"page,omitempty"`
	PerPage uint   `form:"perPage" json:"perPage,omitempty"`
	OrderBy string `form:"order" json:"order,omitempty"`
	Sort    string `form:"sort" json:"sort,omitempty"`
}

func (server *Server) Search(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		logrus.Error("[GetMediaList]: empty userID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var requestBody SearchRequestParams
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
	mediaIDs, err := server.MediaSearch.SearchGetMediaIDsOnly(c.Request.Context(), mediasearch.MediaSearchQuery{
		UserID:  userID,
		OrderBy: mediasearch.ORDER_BY_MEDIA_CREATION_TIME,
		Sort:    mediasearch.SORT_DESCENDING,
		Page:    int64(requestBody.Page),
		PerPage: int64(requestBody.PerPage),
		Query:   requestBody.Query,
	})
	if err != nil {
		logrus.Errorf("[Search] media search failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	mediaList, err := server.Media.GetByMediaIDs(c.Request.Context(), mediaIDs)
	if err != nil {
		logrus.Errorf("[Search] db.GetByMediaIDs failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response := []MediaApiData{}
	for _, media := range mediaList {
		mediaData, err := NewMediaApiData(media)
		if err != nil {
			logrus.Errorf("[Search] NewMediaApiData failed: %w", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		response = append(response, mediaData)
	}
	c.JSON(http.StatusOK, &response)
}

func (requestBody *SearchRequestParams) initDefaultValues() {
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

func (requestBody *SearchRequestParams) validate() error {
	if !utils.Contains(usermediabindings.SUPPORTED_ORDER_BY, requestBody.OrderBy) {
		return fmt.Errorf("invalid orderBy value")
	}
	if requestBody.Sort != usermediabindings.SORT_ASCENDING && requestBody.Sort != usermediabindings.SORT_DESCENDING {
		return fmt.Errorf("invalid sort value")
	}
	// if len(requestBody.Query) == 0 {
	// 	return fmt.Errorf("empty query")
	// }
	return nil
}
