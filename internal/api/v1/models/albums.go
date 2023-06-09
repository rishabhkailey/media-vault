package v1models

import (
	"fmt"
	"time"

	"github.com/rishabhkailey/media-service/internal/services/album"
	"github.com/sirupsen/logrus"
)

type CreateAlbumRequest struct {
	Name         string `form:"name" json:"name" binding:"required"`
	ThumbnailUrl string `form:"thumbnail_url" json:"thumbnail_url"`
}

func (request *CreateAlbumRequest) Validate() error {
	if len(request.Name) == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: empty name")
	}
	return nil
}

type AlbumResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ThumbnailUrl string    `json:"thumbnail_url"`
	MediaCount   int       `json:"media_count"`
}

type GetAlbumsRequest struct {
	OrderBy string `form:"order" json:"order" binding:"required"`
	Sort    string `form:"sort" json:"sort" binding:"required"`
	Page    int64  `form:"page" json:"page" binding:"required"`
	PerPage int64  `form:"perPage" json:"perPage" binding:"required"`
}

func (request *GetAlbumsRequest) Validate() error {
	if _, ok := album.AlbumOrderAttributesMapping[request.OrderBy]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if _, ok := album.AlbumSortKeywordMapping[request.Sort]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
	}
	if request.Page < 0 || request.PerPage < 0 {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: Page or PerPage")
	}
	if request.PerPage > album.MAX_ALBUMS_PER_PAGE_VALUE {
		logrus.Warnf("[MediaSearchQueryValidator] PerPage value exceeded the max supported value")
		request.PerPage = album.MAX_ALBUMS_PER_PAGE_VALUE
	}
	return nil
}

type GetAlbumRequest struct {
	AlbumID uint `uri:"albumID" binding:"required"`
}

func (request *GetAlbumRequest) Validate() error {
	if request.AlbumID == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: invalid album ID")
	}
	return nil
}

type DeleteAlbumRequest struct {
	AlbumID uint `uri:"albumID" binding:"required"`
}

func (request *DeleteAlbumRequest) Validate() error {
	if request.AlbumID == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: invalid album ID")
	}
	return nil
}

type GetAlbumMediaRequest struct {
	AlbumID uint   `uri:"albumID" binding:"required"`
	OrderBy string `form:"order" json:"order" binding:"required"`
	Sort    string `form:"sort" json:"sort" binding:"required"`
	Page    int64  `form:"page" json:"page" binding:"required"`
	PerPage int64  `form:"perPage" json:"perPage" binding:"required"`
}

func (request *GetAlbumMediaRequest) Validate() error {
	if request.AlbumID == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: invalid album ID")
	}
	if _, ok := album.AlbumMediaOrderAttributesMapping[request.OrderBy]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if _, ok := album.AlbumMediaSortKeywordMapping[request.Sort]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
	}
	if request.Page < 0 || request.PerPage < 0 {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: Page or PerPage")
	}
	if request.PerPage > album.MAX_ALBUM_MEDIA_PER_PAGE_VALUE {
		logrus.Warnf("[MediaSearchQueryValidator] PerPage value exceeded the max supported value")
		request.PerPage = album.MAX_ALBUM_MEDIA_PER_PAGE_VALUE
	}
	return nil
}

type AlbumAddMediaRequest struct {
	AlbumID  uint   `uri:"albumID" binding:"required"`
	MediaIDs []uint `json:"media_ids" binding:"required"`
}

func (request *AlbumAddMediaRequest) Validate() error {
	if request.AlbumID == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: invalid album ID")
	}
	if len(request.MediaIDs) == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: empty media_ids")
	}
	return nil
}

type AlbumRemoveMediaRequest struct {
	AlbumID  uint   `uri:"albumID" binding:"required"`
	MediaIDs []uint `json:"media_ids" binding:"required"`
}

func (request *AlbumRemoveMediaRequest) Validate() error {
	if request.AlbumID == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: invalid album ID")
	}
	if len(request.MediaIDs) == 0 {
		return fmt.Errorf("[CreateAlbumRequest.validate]: empty media_ids")
	}
	return nil
}

type AlbumRemoveMediaResponse struct {
	MediaIDs []uint `json:"media_ids"`
}

type AlbumAddMediaResponse struct {
	MediaIDs []uint `json:"media_ids"`
}
