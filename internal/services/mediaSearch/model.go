package mediasearch

import "github.com/rishabhkailey/media-service/internal/services/media"

// todo similar in userMediaBindings
type (
	Order string
	Sort  string
)

const (
	PRIMARY_KEY                        = "media_id"
	USER_ID_KEY                        = "user_id"
	ORDER_BY_UPLOAD_TIME         Order = "uploaded_at"
	ORDER_BY_MEDIA_CREATION_TIME Order = "metadata.timestamp"
	SORT_ASCENDING               Sort  = "asc"
	SORT_DESCENDING              Sort  = "desc"
)

var (
	SearchableAttributes = []string{"metadata"}
	FilterableAttributes = []string{"user_id"}
	SortableAttributes   = []string{"uploaded_at", "metadata.timestamp"}
)

type MeiliSearchMediaMetadata struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

type Model struct {
	MediaID    uint                     `json:"media_id"`
	UserID     string                   `json:"user_id"`
	UploadedAt int64                    `json:"uploaded_at"`
	Metadata   MeiliSearchMediaMetadata `json:"metadata"`
}

type CreateCommand struct {
	Media  media.Model
	UserID string
}

type MediaSearchQuery struct {
	UserID  string
	OrderBy Order
	Sort    Sort
	Page    int64
	PerPage int64
	Query   string
}
