package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Email    string `gorm:";uniqueIndex:users_index_email;"`
	Username string
	Password string
}

type UserModel struct {
	Db *gorm.DB
}

func NewUserModel(db *gorm.DB) (*UserModel, error) {
	err := db.Migrator().AutoMigrate(&Test{})
	if err != nil {
		return nil, err
	}
	return &UserModel{
		Db: db,
	}, nil
}

func (model *UserModel) FindByID(ctx context.Context, id string) (Test Test, err error) {
	err = model.Db.First(&Test, "id = ?", id).Error
	if err != nil {
		return Test, fmt.Errorf("error getting Test details from DB: %w", err)
	}
	return Test, err
}

func (model *UserModel) FindByEmail(ctx context.Context, email string) (Test Test, err error) {
	err = model.Db.First(&Test, "email = ?", email).Error
	if err != nil {
		return Test, fmt.Errorf("error getting Test details from DB: %w", err)
	}
	return Test, err
}

func (model *UserModel) Create(ctx context.Context, email, username, password string) error {
	Test := Test{
		Email:    email,
		Password: password,
		Username: username,
	}
	err := model.Db.WithContext(ctx).Create(&Test).Error
	if err != nil {
		return fmt.Errorf("[Create]: insert failed: %w", err)
	}
	return nil
}
