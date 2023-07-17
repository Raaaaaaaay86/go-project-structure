package mongodb

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type VideoCommentRepository struct {
	Client *mongo.Client
}

func NewVideoCommentRepository(client *mongo.Client) *VideoCommentRepository {
	return &VideoCommentRepository{
		Client: client,
	}
}
func (v VideoCommentRepository) comments() *mongo.Collection {
	return v.Client.Database("video").Collection("comments")
}

func (v VideoCommentRepository) Create(comment *entity.VideoComment) error {
	result, err := v.comments().InsertOne(context.TODO(), comment)
	if err != nil {
		return err
	}
	comment.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (v VideoCommentRepository) FindByVideoId(videoId uint) ([]*entity.VideoComment, error) {
	cursor, err := v.comments().Find(context.TODO(), bson.M{"video_id": videoId})
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

func (v VideoCommentRepository) FindById(id primitive.ObjectID) (*entity.VideoComment, error) {
	result := v.comments().FindOne(context.Background(), bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var comment *entity.VideoComment
	err := result.Decode(&comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (v VideoCommentRepository) DeleteById(id primitive.ObjectID, deleterId uint) (int, error) {
	session, err := v.Client.StartSession()
	if err != nil {
		return 0, err
	}
	defer session.EndSession(context.TODO())

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	deleteCount, err := session.WithTransaction(context.TODO(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		comment, err := v.FindById(id)
		if err != nil {
			return 0, err
		}

		if comment.AuthorId != deleterId {
			return 0, exception.ErrUnauthorized
		}

		deleteResult, err := v.comments().DeleteOne(sessionContext, bson.M{"_id": id})
		if err != nil {
			return nil, err
		}

		return deleteResult.DeletedCount, nil
	}, txnOptions)
	if err != nil {
		return 0, err
	}

	return int(deleteCount.(int64)), nil
}

func (v VideoCommentRepository) ForceDeleteById(id primitive.ObjectID) (int, error) {
	deleteResult, err := v.comments().DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	return int(deleteResult.DeletedCount), nil
}
