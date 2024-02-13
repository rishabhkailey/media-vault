package media

import (
	"time"

	"github.com/rishabhkailey/media-service/internal/store/media"
)

const (
	MAX_PER_PAGE_VALUE int64 = 100
)

var (
	SortKeywordMapping = map[string]media.Sort{
		"asc":        media.Ascending,
		"desc":       media.Descending,
		"ascending":  media.Ascending,
		"descending": media.Descending,
	}
	OrderKeywordMapping = map[string]media.OrderBy{
		"date":        media.Date,
		"uploaded_at": media.UploadedAt,
	}
)

type CreateMediaCommand struct {
	UploadRequestID string
	MetadataID      uint
}

type GetByUploadRequestQuery struct {
	UploadRequestID string
}

type GetByFileNameQuery struct {
	FileName string
}

type GetTypeByFileNameQuery struct {
	FileName string
}

type GetByUserIDQuery struct {
	UserID      string
	OrderBy     string
	Sort        string
	LastMediaID *uint
	LastDate    *time.Time
	PerPage     int64
}

type GetByMediaIDsWithSortQuery struct {
	MediaIDs []uint
	OrderBy  string
	Sort     string
}

type GetByMediaIDQuery struct {
	MediaID uint
}
type GetByMediaIDsQuery struct {
	MediaIDs []uint
}

type UserMediaByIDQuery struct {
	UserID  string
	MediaID uint
}

type DeleteOneCommand struct {
	ID         uint
	UserID     string
	MetadataID uint
}

type DeleteManyCommand struct {
	MediaIDs []uint
	UserID   string
}
