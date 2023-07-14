package repository

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"gorm.io/gorm"
)

type VideoPostRepository interface {
	StartTx() *gorm.DB
	CommitTx(tx *gorm.DB) *gorm.DB
	WithTx(tx *gorm.DB) VideoPostRepository
	ForUpdate() VideoPostRepository
	Create(post *entity.VideoPost) error
	FindById(id uint) (*entity.VideoPost, error)
	Update(post *entity.VideoPost) error
}
