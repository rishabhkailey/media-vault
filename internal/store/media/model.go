package media

type OrderBy string
type Sort string

const (
	Date       OrderBy = "date"
	UploadedAt OrderBy = "created_at"
	Ascending  Sort    = "asc"
	Descending Sort    = "desc"
)
