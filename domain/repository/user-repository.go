package repository

import "github.com/raaaaaaaay86/go-project-structure/domain/entity"

type UserRepository interface {
	WithPreload() UserRepository
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindById(id uint) (*entity.User, error)
}
