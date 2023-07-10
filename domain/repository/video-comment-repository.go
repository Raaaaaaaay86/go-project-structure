package repository

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoCommentRepository interface {
	Create(comment *entity.VideoComment) error
	FindByVideoId(videoId uint) ([]*entity.VideoComment, error)
	FindById(id primitive.ObjectID) (*entity.VideoComment, error)
	DeleteById(id primitive.ObjectID) error
}
