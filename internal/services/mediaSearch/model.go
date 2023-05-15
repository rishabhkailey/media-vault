package mediasearch

import (
	"fmt"

	"github.com/rishabhkailey/media-service/internal/services/media"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
)

const (
	PRIMARY_KEY              = "media_id"
	USER_ID_KEY              = "user_id"
	MAX_PER_PAGE_VALUE int64 = 100
)

var (
	SortKeywordMapping = map[string]string{
		"asc":        "asc",
		"desc":       "desc",
		"ascending":  "asc",
		"descending": "desc",
	}

	OrderAttributesMapping = map[string]string{
		"date":        "metadata.timestamp",
		"uploaded_at": "uploaded_at",
	}

	SearchableAttributes = []string{"metadata"}
	FilterableAttributes = []string{"user_id"}
	SortableAttributes   = []string{"metadata.timestamp", "uploaded_at"}
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
	OrderBy string `form:"order" json:"order,omitempty" binding:"required"`
	Sort    string `form:"sort" json:"sort,omitempty" binding:"required"`
	Page    int64  `form:"page" json:"page,omitempty" binding:"required"`
	PerPage int64  `form:"perPage" json:"perPage,omitempty" binding:"required"`
	Query   string `form:"query" json:"query" binding:"required"`
}

func UserMediaBindingToMeiliSearchMediaIndex(userMediaBindingList []usermediabindings.Model) (meiliSearchMediaList []Model, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[ToMeiliSearchMediaIndex] panic :%v", r)
		}
	}()
	for _, userMediaBinding := range userMediaBindingList {
		meiliSearchMediaList = append(meiliSearchMediaList, Model{
			MediaID:    userMediaBinding.Media.ID,
			UserID:     userMediaBinding.UserID,
			UploadedAt: userMediaBinding.CreatedAt.Unix(),
			Metadata: MediaSearchMetadata{
				Name:      userMediaBinding.Media.Metadata.Name,
				Type:      userMediaBinding.Media.Metadata.Type,
				Timestamp: userMediaBinding.Media.Metadata.Date.Unix(),
				Date:      userMediaBinding.Media.Metadata.Date.Format("Monday January 2 2006 UTC"),
			},
		})
	}
	return
}

func MediaToMeiliSearchMediaIndex(mediaList []media.Model, userID string) (meiliSearchMediaList []Model, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[ToMeiliSearchMediaIndex] panic :%v", r)
		}
	}()
	for _, media := range mediaList {
		meiliSearchMediaList = append(meiliSearchMediaList, Model{
			MediaID:    media.ID,
			UserID:     userID,
			UploadedAt: media.CreatedAt.Unix(),
			Metadata: MediaSearchMetadata{
				Name:      media.Metadata.Name,
				Type:      media.Metadata.Type,
				Timestamp: media.Metadata.Date.Unix(),
				Date:      media.Metadata.Date.Format("Monday January 2 2006 UTC"),
			},
		})
	}
	return
}
