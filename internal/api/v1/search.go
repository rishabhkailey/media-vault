package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1models "github.com/rishabhkailey/media-vault/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-vault/internal/errors"
	"github.com/rishabhkailey/media-vault/internal/services/album"
	albumsearch "github.com/rishabhkailey/media-vault/internal/services/album_search"
	"github.com/rishabhkailey/media-vault/internal/services/media"
	mediasearch "github.com/rishabhkailey/media-vault/internal/services/mediaSearch"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"github.com/sirupsen/logrus"
)

// todo rename to media search. we will also have album search?
func (server *Server) Search(c *gin.Context) {
	userID := c.GetString("user_id")
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
	mediaList, err := server.Media.GetByMediaIDsWithSort(c.Request.Context(), media.GetByMediaIDsWithSortQuery{
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

func (server *Server) SearchAlbum(c *gin.Context) {
	userID := c.GetString("user_id")
	if len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.SearchAlbumRequest
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
	albumIds, err := server.AlbumSearch.GetAlbumIDs(c.Request.Context(), albumsearch.AlbumSearchQuery{
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
	if len(albumIds) == 0 {
		c.JSON(http.StatusOK, v1models.NewAlbumResponse([]storemodels.AlbumModel{}))
		return
	}

	albums, err := server.AlbumService.GetByIdsWithSort(c.Request.Context(), album.GetByIDsWithSortQuery{
		IDs:     albumIds,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
	})
	if err != nil {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search] db.GetByMediaIDs failed: %w", err),
			),
		)
		return
	}
	{
		bytes, err := json.Marshal(v1models.NewAlbumResponse(albums))
		logrus.WithFields(logrus.Fields{
			"err": err,
			"res": string(bytes),
		}).Info("[SearchAlbum] response")
	}
	c.JSON(http.StatusOK, v1models.NewAlbumResponse(albums))
}
