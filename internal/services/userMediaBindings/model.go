package usermediabindings

type CreateCommand struct {
	UserID  string
	MediaID uint
}

type GetByMediaIDQuery struct {
	MediaID uint
}

type CheckFileBelongsToUserQuery struct {
	UserID   string
	FileName string
}

type CheckMediaBelongsToUserQuery struct {
	UserID  string
	MediaID uint
}

type GetUserMediaQuery struct {
	UserID  string
	OrderBy string
	Sort    string
	Offset  int
	Limit   int
}

type DeleteOneCommand struct {
	UserID  string
	MediaID uint
}

type DeleteManyCommand struct {
	UserID   string
	MediaIDs []uint
}
