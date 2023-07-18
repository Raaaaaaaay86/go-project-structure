package repository

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoCommentRepository interface {
	Create(ctx context.Context, comment *entity.VideoComment) error
	FindByVideoId(ctx context.Context, videoId uint) ([]*entity.VideoComment, error)
	FindById(ctx context.Context, id primitive.ObjectID) (*entity.VideoComment, error)
	DeleteById(ctx context.Context, id primitive.ObjectID, deleterId uint) (int, error)
	ForceDeleteById(ctx context.Context, id primitive.ObjectID) (int, error)
}
