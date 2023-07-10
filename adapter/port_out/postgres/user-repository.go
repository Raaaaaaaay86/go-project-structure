package postgres

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u UserRepository) Create(user *entity.User) error {
	tx := u.DB.Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u UserRepository) FindById(id uint) (*entity.User, error) {
	var user entity.User
	tx := u.DB.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (u UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	tx := u.DB.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (u UserRepository) WithPreload() repository.UserRepository {
	u.DB = u.DB.Preload("Roles")
	return u
}
