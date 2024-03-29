package v1models

import (
	"fmt"
	"time"

	"github.com/rishabhkailey/media-vault/internal/services/media"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"github.com/sirupsen/logrus"
)

type GetMediaListRequest struct {
	OrderBy string `form:"order" json:"order" binding:"required"`
	Sort    string `form:"sort" json:"sort" binding:"required"`
	// todo change PerPage to limit or something
	PerPage     int64 `form:"per_page" json:"per_page" binding:"required"`
	LastMediaID *uint `form:"last_media_id" json:"last_media_id"`
	// todo change last_date to last_media_date?
	LastDate *time.Time `form:"last_date" json:"last_date"`
}

func (request *GetMediaListRequest) Validate() error {
	if _, ok := media.OrderKeywordMapping[request.OrderBy]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if _, ok := media.SortKeywordMapping[request.Sort]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
	}
	if request.PerPage < 0 {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: Page or PerPage")
	}
	if request.PerPage > media.MAX_PER_PAGE_VALUE {
		logrus.Warnf("[MediaSearchQueryValidator] PerPage value exceeded the max supported value")
		request.PerPage = media.MAX_PER_PAGE_VALUE
	}
	return nil
}

type GetMediaListResponse []GetMediaResponse
type GetMediaResponse struct {
	Id                   uint      `json:"id"`
	MediaUrl             string    `json:"url"`
	ThumbnailUrl         string    `json:"thumbnail_url"`
	UploadedAt           time.Time `json:"uploaded_at"`
	Name                 string    `json:"name"`
	Date                 time.Time `json:"date"`
	Type                 string    `json:"type"`
	Size                 uint64    `json:"size"`
	Thumbnail            bool      `json:"thumbnail"`
	ThumbnailAspectRatio float32   `json:"thumbnail_aspect_ratio"`
}

func NewGetMediaListResponse(mediaList []storemodels.MediaModel) (result GetMediaListResponse, err error) {
	result = []GetMediaResponse{} // required, if not done then we get null in json
	for _, mediaItem := range mediaList {
		var item GetMediaResponse
		item, err = NewGetMediaResponse(mediaItem)
		if err != nil {
			return
		}
		result = append(result, item)
	}
	return
}

func NewGetMediaResponse(media storemodels.MediaModel) (item GetMediaResponse, err error) {
	item.MediaUrl = parseMediaURL(media.FileName)
	if err != nil {
		return
	}
	item.Id = media.ID
	item.Date = media.Metadata.Date
	item.Name = media.Metadata.Name
	item.Size = media.Metadata.Size
	item.Type = media.Metadata.Type
	item.UploadedAt = media.Metadata.CreatedAt
	item.Thumbnail = media.Metadata.Thumbnail
	item.ThumbnailAspectRatio = media.Metadata.ThumbnailAspectRatio
	if media.Metadata.Thumbnail {
		item.ThumbnailUrl = parseMediaThumbnailUrl(media.FileName)
	}
	return
}

// todo add base url?
func parseMediaThumbnailUrl(fileName string) string {
	return fmt.Sprintf("/v1/file/%s/thumbnail", fileName)
}

func parseMediaURL(fileName string) string {
	return fmt.Sprintf("/v1/file/%s", fileName)
}

type DeleteMediaRequest struct {
	MediaIDs []uint `json:"media_ids" binding:"required"`
}

func (r DeleteMediaRequest) Validate() error {
	if len(r.MediaIDs) == 0 {
		return fmt.Errorf("[DeleteMediaRequest.Validate]: empty media ids slice")
	}
	return nil
}
