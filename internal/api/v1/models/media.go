package v1models

import (
	"fmt"

	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type GetMediaListRequest struct {
	OrderBy media.Order `form:"order" json:"order" binding:"required"`
	Sort    media.Sort  `form:"sort" json:"sort" binding:"required"`
	Page    int64       `form:"page" json:"page" binding:"required"`
	PerPage int64       `form:"perPage" json:"perPage" binding:"required"`
}

func (request *GetMediaListRequest) Validate() error {
	if !utils.Contains(media.SUPPORTED_ORDER_BY, request.OrderBy) {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: OrderBy")
	}
	if !utils.Contains(media.SUPPORTED_SORT_BY, request.Sort) {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: sort")
	}

	if request.Page < 0 || request.PerPage < 0 {
		return fmt.Errorf("[MediaSearchQueryValidator] invalid param: Page or PerPage")
	}
	if request.PerPage > media.MAX_PER_PAGE_VALUE {
		logrus.Warnf("[MediaSearchQueryValidator] PerPage value exceeded the max supported value")
		request.PerPage = media.MAX_PER_PAGE_VALUE
	}
	return nil
}

type GetMediaListResponse []media.GetMediaQueryResultItem
