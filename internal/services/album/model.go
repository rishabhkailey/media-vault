package album

const (
	MAX_ALBUMS_PER_PAGE_VALUE      int64 = 100
	MAX_ALBUM_MEDIA_PER_PAGE_VALUE int64 = 100
)

var (
	AlbumSortKeywordMapping = map[string]string{
		"asc":        "asc",
		"desc":       "desc",
		"ascending":  "asc",
		"descending": "desc",
	}
	AlbumOrderAttributesMapping = map[string]string{
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	AlbumMediaSortKeywordMapping = map[string]string{
		"asc":        "asc",
		"desc":       "desc",
		"ascending":  "asc",
		"descending": "desc",
	}
	AlbumMediaOrderAttributesMapping = map[string]string{
		"added_at": "created_at",
	}
)

type CreateAlbumCmd struct {
	Name         string
	UserID       string
	ThumbnailUrl string
}

type GetUserAlbumsQuery struct {
	UserID  string
	OrderBy string
	Sort    string
	Page    int64
	PerPage int64
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
	UserID  string
	AlbumID uint
	OrderBy string
	Sort    string
	Page    int64
	PerPage int64
}
