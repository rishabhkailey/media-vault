package albumsearch

const (
	ORDER_BY_CREATED_AT  = "created_at"
	ORDER_BY_UPLOADED_AT = "uploaded_at"
	PRIMARY_KEY          = "id"
	USER_ID_KEY          = "user_id"
)

type Index struct {
	Name   string `json:"name"`
	ID     uint   `json:"id"`
	UserID string `json:"user_id"`
}

func GetIndexName() string {
	return "albums"
}

func GetSearchableAttributes() []string {
	return []string{"name"}
}

func GetFilterableAttributes() []string {
	return []string{USER_ID_KEY}
}

func GetSortableAttributes() []string {
	return []string{"created_at", "uploaded_at"}
}
