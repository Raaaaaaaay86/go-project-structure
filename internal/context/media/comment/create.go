package comment

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

var _ validate.Validator = (*CreateCommentCommand)(nil)

type CreateCommentCommand struct {
	Comment  string `json:"comment"`
	VideoId  uint   `json:"videoId"`
	AuthorId uint   `json:"-"`
}

func (c CreateCommentCommand) Validate() error {
	if c.Comment == "" {
		return exception.NewInvalidInputError("comment").ShouldNotEmpty()
	}
	return nil
}

type CreateCommentResponse struct {
	CommendId string `json:"commend_id"`
}

type ICreateCommentUseCase interface {
	Execute(ctx context.Context, cmd CreateCommentCommand) (*CreateCommentResponse, error)
}

var _ ICreateCommentUseCase = (*CreateCommentUseCase)(nil)

type CreateCommentUseCase struct {
	VideoCommentRepo    repository.VideoCommentRepository
	UserRepository      repository.UserRepository
	VideoPostRepository repository.VideoPostRepository
	TracerProvider      *trace.TracerProvider
}

func NewCreateCommentUseCase(tracerProvider *trace.TracerProvider, videoCommentRepo repository.VideoCommentRepository, userRepository repository.UserRepository, postRepository repository.VideoPostRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		VideoCommentRepo:    videoCommentRepo,
		UserRepository:      userRepository,
		VideoPostRepository: postRepository,
		TracerProvider:      tracerProvider,
	}
}

func (c CreateCommentUseCase) Execute(ctx context.Context, cmd CreateCommentCommand) (*CreateCommentResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(c.TracerProvider, ctx, pkg, "CreateCommentUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	author, err := c.UserRepository.FindById(newCtx, cmd.AuthorId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	video, err := c.VideoPostRepository.FindById(newCtx, cmd.VideoId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	now := time.Now()
	comment := entity.NewVideoComment(*video, *author, cmd.Comment, now, now)

	err = c.VideoCommentRepo.Create(newCtx, comment)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &CreateCommentResponse{CommendId: comment.Id.Hex()}, nil
}
