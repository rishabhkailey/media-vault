package storemodels

import (
	"github.com/rishabhkailey/media-vault/internal/constants"
	"gorm.io/gorm"
)

type MediaModel struct {
	gorm.Model
	FileName        string `gorm:"index,unique"`
	UploadRequestID string `gorm:"index"`
	UploadRequest   UploadRequestsModel
	Metadata        MediaMetadataModel
	MetadataID      uint
}

func (MediaModel) TableName() string {
	return constants.MEDIA_TABLE
}
