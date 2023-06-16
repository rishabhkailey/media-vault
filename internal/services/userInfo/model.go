package userinfo

type GetUserInfoQuery struct {
	UserID string
}

type CreateUserInfoCmd struct {
	UserID                string
	PreferedTimeZone      string
	EncryptionKeyChecksum string
}
