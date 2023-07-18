package repository

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"gorm.io/gorm"
)

type VideoPostRepository interface {
	StartTx() *gorm.DB
	CommitTx(tx *gorm.DB) *gorm.DB
	WithTx(tx *gorm.DB) VideoPostRepository
	ForUpdate() VideoPostRepository
	Create(ctx context.Context, post *entity.VideoPost) error
	FindById(ctx context.Context, id uint) (*entity.VideoPost, error)
	Update(ctx context.Context, post *entity.VideoPost) error
}
