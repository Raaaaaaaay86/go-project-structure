package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repository.VideoPostRepository = (*VideoPostRepository)(nil)

type VideoPostRepository struct {
	DB             *gorm.DB
	GraphDB        neo4j.DriverWithContext
	TracerProvider *trace.TracerProvider
}

func NewVideoPostRepository(tracerProvider *trace.TracerProvider, DB *gorm.DB, graphDB neo4j.DriverWithContext) *VideoPostRepository {
	return &VideoPostRepository{
		DB:             DB,
		GraphDB:        graphDB,
		TracerProvider: tracerProvider,
	}
}

func (v VideoPostRepository) Like(ctx context.Context, videoId uint, userId uint) error {
	session := v.GraphDB.NewSession(ctx, neo4j.SessionConfig{})

	_, err := neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		params := map[string]any{
			"videoId": videoId,
			"userId":  userId,
		}

		_, err := tx.Run(ctx, "MERGE (u:User{id:$userId}) MERGE (v:Video{id:$videoId}) MERGE (u)-[:LIKE]->(v)", params)
		if err != nil {
			return nil, err
		}

		return 1, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (v VideoPostRepository) UnLike(ctx context.Context, videoId uint, userId uint) error {
	session := v.GraphDB.NewSession(ctx, neo4j.SessionConfig{})

	_, err := neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		params := map[string]any{
			"videoId": videoId,
			"userId":  userId,
		}

		_, err := tx.Run(ctx, "MATCH (u:User{id:$userId})-[r:LIKE]->(v:Video{id:$videoId}) DELETE r", params)
		if err != nil {
			return nil, err
		}

		return 1, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (v VideoPostRepository) ForUpdate() repository.VideoPostRepository {
	v.DB = v.DB.Clauses(clause.Locking{Strength: "UPDATE"})
	return v
}

func (v VideoPostRepository) Update(ctx context.Context, post *entity.VideoPost) error {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.Update")
	defer span.End()

	err := v.DB.WithContext(newCtx).Transaction(func(tx *gorm.DB) error {
		video, err := v.FindById(newCtx, post.Id)
		if err != nil {
			return err
		}

		result := v.DB.WithContext(newCtx).Updates(post).Where("id = ?", video.Id)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (v VideoPostRepository) Create(ctx context.Context, post *entity.VideoPost) error {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.Create")
	defer span.End()

	tx := v.DB.WithContext(newCtx).Create(post)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (v VideoPostRepository) FindById(ctx context.Context, id uint) (*entity.VideoPost, error) {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.FindById")
	defer span.End()

	var post entity.VideoPost
	tx := v.DB.WithContext(newCtx).Where("id = ?", id).First(&post)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &post, nil
}
