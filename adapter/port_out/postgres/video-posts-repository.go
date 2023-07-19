package postgres

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoPostRepository struct {
	DB             *gorm.DB
	TracerProvider *trace.TracerProvider
}

func NewVideoPostRepository(tracerProvider *trace.TracerProvider, DB *gorm.DB) *VideoPostRepository {
	return &VideoPostRepository{
		DB:             DB,
		TracerProvider: tracerProvider,
	}
}

func (v VideoPostRepository) ForUpdate() repository.VideoPostRepository {
	v.DB = v.DB.Clauses(clause.Locking{Strength: "UPDATE"})
	return v
}

func (v VideoPostRepository) Update(ctx context.Context, post *entity.VideoPost) error {
	newCtx, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.Update")
	defer span.End()

	err := v.DB.Transaction(func(tx *gorm.DB) error {
		video, err := v.FindById(newCtx, post.Id)
		if err != nil {
			return err
		}

		result := v.DB.Updates(post).Where("id = ?", video.Id)
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
	_, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.Create")
	defer span.End()

	tx := v.DB.Create(post)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (v VideoPostRepository) FindById(ctx context.Context, id uint) (*entity.VideoPost, error) {
	_, span := tracing.RepositorySpanFactory(v.TracerProvider, ctx, pkg, "VideoPostRepository.FindById")
	defer span.End()

	var post entity.VideoPost
	tx := v.DB.Where("id = ?", id).First(&post)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &post, nil
}
