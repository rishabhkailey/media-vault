package db

import (
	dbservices "github.com/rishabhkailey/media-service/internal/db/services"
	"gorm.io/gorm"
)

type Services struct {
	Users dbservices.UserModel
}

func NewServices(db *gorm.DB) (*Services, error) {
	userModel, err := dbservices.NewUserModel(db)
	if err != nil {
		return nil, err
	}

	return &Services{
		Users: *userModel,
	}, nil
}
