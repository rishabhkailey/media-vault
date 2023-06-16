package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	v1models "github.com/rishabhkailey/media-service/internal/api/v1/models"
	internalErrors "github.com/rishabhkailey/media-service/internal/errors"
	"github.com/rishabhkailey/media-service/internal/services/album"
)

func (server *Server) CreateAlbum(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.CreateAlbumRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[CreateAlbum] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[CreateAlbum] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	album, err := server.AlbumService.Create(c.Request.Context(), album.CreateAlbumCmd{
		Name:         requestBody.Name,
		ThumbnailUrl: requestBody.ThumbnailUrl,
		UserID:       userID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	response := v1models.AlbumResponse{
		ID:           album.ID,
		Name:         album.Name,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
		ThumbnailUrl: album.ThumbnailUrl,
		MediaCount:   album.MediaCount,
	}
	c.JSON(http.StatusOK, &response)
}

func (server *Server) GetAlbums(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.GetAlbumsRequest
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	albums, err := server.AlbumService.GetUserAlbums(c.Request.Context(), album.GetUserAlbumsQuery{
		UserID:  userID,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		c.Error(err)
		return
	}
	albumsResponse := []v1models.AlbumResponse{}
	for _, album := range albums {
		albumsResponse = append(albumsResponse, v1models.AlbumResponse{
			ID:           album.ID,
			Name:         album.Name,
			CreatedAt:    album.CreatedAt,
			UpdatedAt:    album.UpdatedAt,
			ThumbnailUrl: album.ThumbnailUrl,
			MediaCount:   album.MediaCount,
		})
	}
	c.JSON(http.StatusOK, &albumsResponse)
}

func (server *Server) DeleteAlbum(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.DeleteAlbumRequest
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	err := server.AlbumService.DeleteAlbum(c.Request.Context(), album.DeleteAlbumCmd{
		AlbumID: requestBody.AlbumID,
		UserID:  userID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func (server *Server) GetAlbum(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.GetAlbumRequest
	if err := c.BindUri(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	album, err := server.AlbumService.GetUserAlbum(c.Request.Context(), album.GetUserAlbumQuery{
		UserID:  userID,
		AlbumID: requestBody.AlbumID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.AlbumResponse{
		ID:           album.ID,
		Name:         album.Name,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
		ThumbnailUrl: album.ThumbnailUrl,
		MediaCount:   album.MediaCount,
	})
}

func (server *Server) PathchAlbum(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.PatchAlbumRequest
	if err := c.BindUri(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := c.Bind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	album, err := server.AlbumService.UpdateAlbum(c.Request.Context(), album.UpdateAlbumCmd{
		AlbumID:      requestBody.ID,
		Name:         requestBody.Name,
		ThumbnailUrl: requestBody.ThumbnailUrl,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.AlbumResponse{
		ID:           album.ID,
		Name:         album.Name,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
		ThumbnailUrl: album.ThumbnailUrl,
		MediaCount:   album.MediaCount,
	})
}

func (server *Server) GetAlubmMedia(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.GetAlbumMediaRequest
	_ = c.ShouldBindUri(&requestBody)
	if err := c.ShouldBind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	response, err := server.AlbumService.GetAlbumMedia(c.Request.Context(), album.GetAlbumMediaQuery{
		UserID:  userID,
		AlbumID: requestBody.AlbumID,
		OrderBy: requestBody.OrderBy,
		Sort:    requestBody.Sort,
		Page:    requestBody.Page,
		PerPage: requestBody.PerPage,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (server *Server) AlbumAddMedia(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.AlbumAddMediaRequest
	_ = c.ShouldBindUri(&requestBody)
	if err := c.ShouldBind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	addedMediaIDs, err := server.AlbumService.AddMedia(c.Request.Context(), album.AddMediaQuery{
		AlbumID:  requestBody.AlbumID,
		UserID:   userID,
		MediaIDs: requestBody.MediaIDs,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.AlbumAddMediaResponse{
		MediaIDs: addedMediaIDs,
	})
}

func (server *Server) RemoveAlbumMedia(c *gin.Context) {
	userID, ok := c.Keys["userID"].(string)
	if !ok || len(userID) == 0 {
		c.Error(
			internalErrors.NewInternalServerError(
				fmt.Errorf("[Search]: empty userID"),
			),
		)
		return
	}
	var requestBody v1models.AlbumRemoveMediaRequest
	_ = c.ShouldBindUri(&requestBody)
	if err := c.ShouldBind(&requestBody); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
	}
	if err := requestBody.Validate(); err != nil {
		c.Error(
			internalErrors.NewBadRequestError(
				fmt.Errorf("[GetAlbums] invalid request: %w", err),
				"bad request",
			),
		)
		return
	}
	removedMediaIDs, err := server.AlbumService.RemoveAlbumMedia(c.Request.Context(), album.RemoveMediaCmd{
		AlbumID:  requestBody.AlbumID,
		UserID:   userID,
		MediaIDs: requestBody.MediaIDs,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &v1models.AlbumRemoveMediaResponse{
		MediaIDs: removedMediaIDs,
	})
}
