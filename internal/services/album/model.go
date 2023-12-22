package album

import (
	"time"

	"github.com/rishabhkailey/media-service/internal/store/album"
)

const (
	MAX_ALBUMS_PER_PAGE_VALUE      int   = 100
	MAX_ALBUM_MEDIA_PER_PAGE_VALUE int64 = 100
)

var (
	AlbumSortKeywordMapping = map[string]album.Sort{
		"asc":        album.Ascending,
		"desc":       album.Descending,
		"ascending":  album.Ascending,
		"descending": album.Descending,
	}
	// user friendly order value to store attribute (column in case of SQL)
	AlbumOrderAttributesMapping = map[string]album.AlbumOrderBy{
		"created_at": album.AlbumOrderByDate,
		"updated_at": album.AlbumOrderByUpdatedAt,
	}
	AlbumMediaSortKeywordMapping = map[string]string{
		"asc":        "asc",
		"desc":       "desc",
		"ascending":  "asc",
		"descending": "desc",
	}
	// user friendly order value to store attribute (column in case of SQL)
	AlbumMediaOrderAttributesMapping = map[string]album.AlbumMediaOrderBy{
		"added_at":    album.AlbumMediaOrderByAddedDate,
		"uploaded_at": album.AlbumMediaOrderByUploadedDate,
		"date":        album.AlbumMediaOrderByMediaDate,
	}
)

type CreateAlbumCmd struct {
	Name         string
	UserID       string
	ThumbnailUrl string
}

type GetUserAlbumsQuery struct {
	UserID        string
	OrderBy       album.AlbumOrderBy
	Sort          album.Sort
	PerPage       int
	LastAlbumID   *uint
	LastAlbumDate *time.Time
}

type GetUserAlbumQuery struct {
	UserID  string
	AlbumID uint
}

type AddMediaQuery struct {
	AlbumID  uint
	UserID   string
	MediaIDs []uint
}

type RemoveMediaCmd struct {
	AlbumID  uint
	UserID   string
	MediaIDs []uint
}

type DeleteAlbumCmd struct {
	AlbumID uint
	UserID  string
}

type GetAlbumMediaQuery struct {
	UserID      string
	AlbumID     uint
	OrderBy     album.AlbumMediaOrderBy
	Sort        album.Sort
	PerPage     int64
	LastMediaID *uint
	LastDate    *time.Time
}

type UpdateAlbumCmd struct {
	AlbumID      uint
	Name         *string
	ThumbnailUrl *string
}
