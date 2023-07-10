package postgres

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"gorm.io/gorm"
)

type VideoPostRepository struct {
	DB *gorm.DB
}

func NewVideoPostRepository(DB *gorm.DB) *VideoPostRepository {
	return &VideoPostRepository{DB: DB}
}

func (v VideoPostRepository) Create(post *entity.VideoPost) error {
	tx := v.DB.Create(post)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (v VideoPostRepository) FindById(id uint) (*entity.VideoPost, error) {
	var post entity.VideoPost
	tx := v.DB.Where("id = ?", id).First(&post)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &post, nil
}
