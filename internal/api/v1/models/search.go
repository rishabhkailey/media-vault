package v1models

import (
	"fmt"
	"slices"

	"github.com/rishabhkailey/media-vault/internal/services/album"
	mediasearch "github.com/rishabhkailey/media-vault/internal/services/mediaSearch"
	"github.com/sirupsen/logrus"
)

type SearchRequest struct {
	OrderBy string `form:"order" json:"order,omitempty" binding:"required"`
	Sort    string `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64  `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64  `form:"per_page" json:"per_page,omitempty" binding:"required"`
	Query   string `form:"query" json:"query" binding:"required"`
}

func (request *SearchRequest) Validate() error {
	if _, ok := mediasearch.OrderAttributesMapping[request.OrderBy]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if _, ok := mediasearch.SortKeywordMapping[request.Sort]; !ok {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
	}
	if request.Page < 0 || request.PerPage < 0 {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: Page or PerPage")
	}
	if request.PerPage > mediasearch.MAX_PER_PAGE_VALUE {
		logrus.Warnf("[MediaSearchQueryValidator] PerPage value exceeded the max supported value")
		request.PerPage = mediasearch.MAX_PER_PAGE_VALUE
	}
	return nil
}

type SearchResponse []GetMediaResponse

type SearchAlbumRequest struct {
	OrderBy string `form:"order" json:"order,omitempty" binding:"required"`
	Sort    string `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64  `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64  `form:"per_page" json:"per_page,omitempty" binding:"required"`
	Query   string `form:"query" json:"query" binding:"required"`
}

func (request *SearchAlbumRequest) Validate() error {
	if !slices.Contains(album.GetValidOrderByValues(), request.OrderBy) {
		return fmt.Errorf(`[MediaSearchQueryValidator] invalid param: "order"`)
	}
	if !slices.Contains(album.GetValidSortValues(), request.Sort) {
		return fmt.Errorf(`[MediaSearchQueryValidator] invalid param: "sort"`)
	}
	if request.Page < 0 || request.PerPage < 0 {
		return fmt.Errorf(`[MediaSearchQueryValidator] invalid param: "page" or "per_page"`)
	}
	if request.PerPage > int64(album.MAX_ALBUMS_PER_PAGE_VALUE) {
		logrus.WithFields(logrus.Fields{
			"per_page":       request.PerPage,
			"per_page_limit": int64(album.MAX_ALBUMS_PER_PAGE_VALUE),
		}).Warnf(`[MediaSearchQueryValidator] "per_page" value exceeded the max supported value`)
		request.PerPage = mediasearch.MAX_PER_PAGE_VALUE
	}
	return nil
}

type SearchAlbumResponse []AlbumResponse
