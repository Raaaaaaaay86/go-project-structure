package mongodb

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoCommentRepository struct {
	Collection *mongo.Collection
}

func NewVideoCommentRepository(collection *mongo.Collection) *VideoCommentRepository {
	return &VideoCommentRepository{Collection: collection}
}

func (v VideoCommentRepository) Create(comment *entity.VideoComment) error {
	result, err := v.Collection.InsertOne(context.TODO(), comment)
	if err != nil {
		return err
	}
	comment.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (v VideoCommentRepository) FindByVideoId(videoId uint) ([]*entity.VideoComment, error) {
	cursor, err := v.Collection.Find(context.TODO(), bson.M{"video_id": videoId})
	if err != nil {
		return nil, err
	}

	var comments []*entity.VideoComment
	for cursor.Next(context.Background()) {
		var comment *entity.VideoComment
		err := cursor.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
