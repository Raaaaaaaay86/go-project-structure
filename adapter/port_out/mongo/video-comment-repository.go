package mongodb

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.opentelemetry.io/otel/sdk/trace"
)

type VideoCommentRepository struct {
	Client         *mongo.Client
	TracerProvider *trace.TracerProvider
}

func NewVideoCommentRepository(tracerProvider *trace.TracerProvider, client *mongo.Client) *VideoCommentRepository {
	return &VideoCommentRepository{
		Client:         client,
		TracerProvider: tracerProvider,
	}
}
func (v VideoCommentRepository) comments() *mongo.Collection {
	return v.Client.Database("video").Collection("comments")
}

func (v VideoCommentRepository) Create(ctx context.Context, comment *entity.VideoComment) error {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoCommentRepository.Create")
	defer span.End()

	result, err := v.comments().InsertOne(newCtx, comment)
	if err != nil {
		return err
	}

	comment.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (v VideoCommentRepository) FindByVideoId(ctx context.Context, videoId uint) ([]*entity.VideoComment, error) {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoCommentRepository.FindByVideoId")
	defer span.End()

	cursor, err := v.comments().Find(newCtx, bson.M{"video_id": videoId})
	if err != nil {
		return nil, err
	}

	var comments []*entity.VideoComment
	for cursor.Next(newCtx) {
		var comment *entity.VideoComment
		err := cursor.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (v VideoCommentRepository) FindById(ctx context.Context, id primitive.ObjectID) (*entity.VideoComment, error) {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoCommentRepository.FindById")
	defer span.End()

	result := v.comments().FindOne(newCtx, bson.M{"_id": id})
	if result.Err() != nil {
		span.RecordError(result.Err())
		return nil, result.Err()
	}

	var comment *entity.VideoComment
	err := result.Decode(&comment)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return comment, nil
}

func (v VideoCommentRepository) DeleteById(ctx context.Context, id primitive.ObjectID, deleterId uint) (int, error) {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoCommentRepository.DeleteById")
	defer span.End()

	session, err := v.Client.StartSession()
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	defer session.EndSession(newCtx)

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	deleteCount, err := session.WithTransaction(newCtx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		comment, err := v.FindById(newCtx, id)
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
		span.RecordError(err)
		return 0, err
	}

	return int(deleteCount.(int64)), nil
}

func (v VideoCommentRepository) ForceDeleteById(ctx context.Context, id primitive.ObjectID) (int, error) {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoCommentRepository.ForceDeleteById")
	defer span.End()

	deleteResult, err := v.comments().DeleteOne(newCtx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	return int(deleteResult.DeletedCount), nil
}
