package userinfo

type UserInfo struct {
	ID                    string `gorm:"primaryKey"`
	PreferedTimeZone      string
	EncryptionKeyChecksum string
	StorageUsage          int64
}
