package postgres

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB             *gorm.DB
	TracerProvider *trace.TracerProvider
}

func NewUserRepository(tracerProvider *trace.TracerProvider, db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB:             db,
		TracerProvider: tracerProvider,
	}
}

func (u UserRepository) Create(ctx context.Context, user *entity.User) error {
	_, span := tracing.RepositorySpanFactory(u.TracerProvider, ctx, pkg, "UserRepository.Create")
	defer span.End()

	tx := u.DB.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u UserRepository) FindById(ctx context.Context, id uint) (*entity.User, error) {
	_, span := tracing.RepositorySpanFactory(u.TracerProvider, ctx, pkg, "UserRepository.FindById")
	defer span.End()

	var user entity.User
	tx := u.DB.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func (u UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	_, span := tracing.RepositorySpanFactory(u.TracerProvider, ctx, pkg, "UserRepository.FindByUsername")
	defer span.End()

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
