package comment

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
)

type FindByVideoQuery struct {
	VideoId uint `form:"videoId"`
}

func (f FindByVideoQuery) Validate() error {
	if f.VideoId == 0 {
		return exception.ErrEmptyInput
	}
	return nil
}

type FindByVideoResponse struct {
	Comments []*entity.VideoComment `json:"comments"`
}

type IFindByVideoCQRS interface {
	Execute(ctx context.Context, query FindByVideoQuery) (*FindByVideoResponse, error)
}

type FindByVideoCQRS struct {
	VideoCommentRepository repository.VideoCommentRepository
	TracerProvider         *trace.TracerProvider
}

func NewFindByVideoUseCase(tracerProvider *trace.TracerProvider, videoCommentRepository repository.VideoCommentRepository) *FindByVideoCQRS {
	return &FindByVideoCQRS{
		VideoCommentRepository: videoCommentRepository,
		TracerProvider:         tracerProvider,
	}
}

func (f FindByVideoCQRS) Execute(ctx context.Context, query FindByVideoQuery) (*FindByVideoResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(f.TracerProvider, ctx, pkg, "FindByVideoCQRS.Execute")
	defer span.End()

	err := validate.Do(query)
	if err != nil {
		return nil, err
	}

	comments, err := f.VideoCommentRepository.FindByVideoId(newCtx, query.VideoId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &FindByVideoResponse{comments}, nil
}
