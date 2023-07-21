package repository

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindById(ctx context.Context, id uint) (*entity.User, error)
}
