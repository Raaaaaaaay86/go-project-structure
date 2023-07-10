package repository

import "github.com/raaaaaaaay86/go-project-structure/domain/entity"

type VideoPostRepository interface {
	Create(post *entity.VideoPost) error
	FindById(id uint) (*entity.VideoPost, error)
}
