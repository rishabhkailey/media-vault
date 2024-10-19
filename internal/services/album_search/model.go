package albumsearch

type AlbumSearchQuery struct {
	UserID  string
	OrderBy string
	Sort    string
	Page    int64
	PerPage int64
	Query   string
}
