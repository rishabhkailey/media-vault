package v1models

import (
	"fmt"
	"time"

	"github.com/rishabhkailey/media-service/internal/services/album"
	albumStore "github.com/rishabhkailey/media-service/internal/store/album"
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
	OrderBy     string `form:"order" json:"order" binding:"required"`
	Sort        string `form:"sort" json:"sort" binding:"required"`
	PerPage     int    `form:"per_page" json:"per_page" binding:"required"`
	LastAlbumID *uint  `form:"last_album_id" json:"last_album_id"`
}

func (request *GetAlbumsRequest) Validate() error {
	if _, ok := album.AlbumOrderAttributesMapping[request.OrderBy]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if _, ok := album.AlbumSortKeywordMapping[request.Sort]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
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

type PatchAlbumRequest struct {
	ID           uint    `uri:"albumID" binding:"required"`
	Name         *string `json:"name"`
	ThumbnailUrl *string `json:"thumbnail_url"`
}

func (request *PatchAlbumRequest) Validate() error {
	if request.ID == 0 {
		return fmt.Errorf("[PatchAlbumRequest.validate]: invalid album ID")
	}
	if (request.Name == nil || len(*request.Name) == 0) && (request.ThumbnailUrl == nil || len(*request.ThumbnailUrl) == 0) {
		return fmt.Errorf("[PatchAlbumRequest.validate]: atleast one update feild required")
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
	AlbumID     uint   `uri:"albumID" binding:"required"`
	OrderBy     string `form:"order" json:"order" binding:"required"`
	Sort        string `form:"sort" json:"sort" binding:"required"`
	PerPage     int64  `form:"per_page" json:"per_page" binding:"required"`
	LastMediaID *uint  `form:"last_media_id" json:"last_media_id"`
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

type GetAlbumMediaListResponse []SingleAlbumMediaResponse
type SingleAlbumMediaResponse struct {
	GetMediaResponse
	AddedAt time.Time `json:"added_at"`
}

func NewAlbumMediaListResponse(albumMediaBindings []albumStore.AlbumMediaBindings) (GetAlbumMediaListResponse, error) {
	var response []SingleAlbumMediaResponse

	for _, albumMediaBinding := range albumMediaBindings {
		mediaResponse, err := NewGetMediaResponse(albumMediaBinding.Media)
		if err != nil {
			return response, err
		}
		response = append(response, SingleAlbumMediaResponse{
			AddedAt:          albumMediaBinding.CreatedAt,
			GetMediaResponse: mediaResponse,
		})
	}

	return response, nil
}
