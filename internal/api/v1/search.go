package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/rishabhkailey/media-service/internal/services/media"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
)

// todo rename to media search. we will also have album search?
func (server *Server) Search(c *gin.Context) {
	userID := c.GetString("userID")
	if len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.SearchRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[Search] invalid request: %v", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[Search] invalid request: %v", err),
				"bad request",
			),
		)
		return
	}
	mediaIDs, err := server.MediaSearch.SearchGetMediaIDsOnly(c.Request.Context(), mediasearch.MediaSearchQuery{
		UserID:  userID,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
		Query:   requestBody.Query,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search] media search failed: %w", err),
			),
		)
		return
	}
	mediaList, err := server.Media.GetByMediaIDs(c.Request.Context(), media.GetByMediaIDsQuery{
		MediaIDs: mediaIDs,
		OrderBy:  requestBody.OrderBy,
		Sort:     requestBody.Sort,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search] db.GetByMediaIDs failed: %w", err),
			),
		)
		return
	}

	response, err := v1models.NewGetMediaListResponse(mediaList)
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search] NewGetMediaListResponse failed: %w", err),
			),
		)
		return
	}

	c.JSON(http.StatusOK, response)
}
