package v1models

type GetUserInfoResponse struct {
	ID                    string `json:"id"`
	PreferedTimeZone      string `json:"prefered_timezone"`
	StorageUsage          int64  `json:"storage_usage"`
	EncryptionKeyChecksum string `json:"encryption_key_checksum"`
}

type PostUserInfoRequest struct {
	PreferedTimeZone      string `json:"prefered_timezone" binding:"required"`
	EncryptionKeyChecksum string `json:"encryption_key_checksum" binding:"required"`
}

func (request *PostUserInfoRequest) Validate() error {
	// no checks binding required
	return nil
}

type PostUserInfoResponse struct {
	ID                    string `json:"id"`
	PreferedTimeZone      string `json:"prefered_timezone"`
	StorageUsage          int64  `json:"storage_usage"`
	EncryptionKeyChecksum string `json:"encryption_key_checksum"`
}
