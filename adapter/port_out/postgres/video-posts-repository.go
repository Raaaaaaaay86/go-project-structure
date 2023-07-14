package postgres

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoPostRepository struct {
	DB *gorm.DB
}

func NewVideoPostRepository(DB *gorm.DB) *VideoPostRepository {
	return &VideoPostRepository{DB: DB}
}

func (v VideoPostRepository) StartTx() *gorm.DB {
	return v.DB.Begin()
}

func (v VideoPostRepository) CommitTx(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (v VideoPostRepository) WithTx(tx *gorm.DB) repository.VideoPostRepository {
	return &VideoPostRepository{
		DB: tx,
	}
}

func (v VideoPostRepository) ForUpdate() repository.VideoPostRepository {
	return &VideoPostRepository{
		DB: v.DB.Clauses(clause.Locking{Strength: "UPDATE"}),
	}
}

func (v VideoPostRepository) Update(post *entity.VideoPost) error {
	tx := v.DB.Updates(post)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
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
