package v1models

import (
	"fmt"

	"github.com/rishabhkailey/media-service/internal/services/media"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type SearchRequest struct {
	OrderBy mediasearch.Order `form:"order" json:"order,omitempty" binding:"required"`
	Sort    mediasearch.Sort  `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64             `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64             `form:"perPage" json:"perPage,omitempty" binding:"required"`
	Query   string            `form:"query" json:"query" binding:"required"`
}

func (request *SearchRequest) Validate() error {
	if !utils.Contains(mediasearch.SUPPORTED_ORDER_BY, request.OrderBy) {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if !utils.Contains(mediasearch.SUPPORTED_SORT, request.Sort) {
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

type SearchResponse []media.GetMediaQueryResultItem
