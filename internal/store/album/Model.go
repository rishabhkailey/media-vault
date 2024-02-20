package album

// AlbumOrderBy represent the attributes of Album table that can be used to sort albums
type AlbumOrderBy string

// AlbumOrderBy represent the attributes of AlbumMediaBindings store that can be used to sort media of a album
type AlbumMediaOrderBy string
type Sort string

const (
	AlbumOrderByDate              AlbumOrderBy      = "album.created_at"
	AlbumOrderByUpdatedAt         AlbumOrderBy      = "album.updated_at"
	AlbumMediaOrderByAddedDate    AlbumMediaOrderBy = "album_media_bindings.created_at"
	AlbumMediaOrderByUploadedDate AlbumMediaOrderBy = "media.created_at"
	AlbumMediaOrderByMediaDate    AlbumMediaOrderBy = "media_metadata.date"
	Ascending                     Sort              = "asc"
	Descending                    Sort              = "desc"
)
