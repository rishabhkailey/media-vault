package mediasearch

import (
	"github.com/rishabhkailey/media-service/internal/services/media"
	"github.com/rishabhkailey/media-service/internal/utils"
)

// todo similar in userMediaBindings
type (
	Order string
	Sort  string
)

const (
	PRIMARY_KEY                        = "media_id"
	USER_ID_KEY                        = "user_id"
	ORDER_BY_UPLOAD_TIME         Order = "uploaded_at"
	ORDER_BY_MEDIA_CREATION_TIME Order = "date"
	SORT_ASCENDING               Sort  = "asc"
	SORT_DESCENDING              Sort  = "desc"
	MAX_PER_PAGE_VALUE           int64 = 100
)

var (
	SUPPORTED_ORDER_BY = []Order{ORDER_BY_UPLOAD_TIME, ORDER_BY_MEDIA_CREATION_TIME}
	SUPPORTED_SORT     = []Sort{SORT_ASCENDING, SORT_DESCENDING}

	OrderAttributesMapping = map[Order]string{
		ORDER_BY_MEDIA_CREATION_TIME: "metadata.timestamp",
		ORDER_BY_UPLOAD_TIME:         "uploaded_at",
	}

	SearchableAttributes = []string{"metadata"}
	FilterableAttributes = []string{"user_id"}
	SortableAttributes   = utils.SliceMap(SUPPORTED_ORDER_BY, func(order Order) string {
		return OrderAttributesMapping[order]
	})
)

type MediaSearchMetadata struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

type Model struct {
	MediaID    uint                `json:"media_id"`
	UserID     string              `json:"user_id"`
	UploadedAt int64               `json:"uploaded_at"`
	Metadata   MediaSearchMetadata `json:"metadata"`
}

type CreateCommand struct {
	Media  media.Model
	UserID string
}

type MediaSearchQuery struct {
	UserID  string
	OrderBy Order  `form:"order" json:"order,omitempty" binding:"required"`
	Sort    Sort   `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64  `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64  `form:"perPage" json:"perPage,omitempty" binding:"required"`
	Query   string `form:"query" json:"query" binding:"required"`
}
