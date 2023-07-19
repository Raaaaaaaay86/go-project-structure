package repository

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
)

type VideoPostRepository interface {
	Create(ctx context.Context, post *entity.VideoPost) error
	FindById(ctx context.Context, id uint) (*entity.VideoPost, error)
	Update(ctx context.Context, post *entity.VideoPost) error
}
