package video

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ validate.Validator = (*CreateVideoCommand)(nil)

type CreateVideoCommand struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	VideoUUID   string `json:"uuid,omitempty"`
	AuthorId    uint   `json:"authorId,omitempty"`
}

func (c CreateVideoCommand) Validate() error {
	if c.AuthorId == 0 {
		return exception.ErrEmptyInput
	}

	fields := []string{c.Title, c.VideoUUID}
	for _, field := range fields {
		if len(field) == 0 {
			return exception.ErrEmptyInput
		}
	}
	return nil
}

type CreateVideoResponse struct {
}

var _ ICreateVideoUseCase = (*CreateVideoUseCase)(nil)

type ICreateVideoUseCase interface {
	Execute(ctx context.Context, cmd CreateVideoCommand) (*CreateVideoResponse, error)
}

type CreateVideoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	UserRepository      repository.UserRepository
	TracerProvider      *trace.TracerProvider
}

func NewCreateVideoUseCase(tracerProvider *trace.TracerProvider, videoPostRepository repository.VideoPostRepository, userRepository repository.UserRepository) *CreateVideoUseCase {
	return &CreateVideoUseCase{
		VideoPostRepository: videoPostRepository,
		UserRepository:      userRepository,
		TracerProvider:      tracerProvider,
	}
}

func (v CreateVideoUseCase) Execute(ctx context.Context, cmd CreateVideoCommand) (*CreateVideoResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(v.TracerProvider, ctx, pkg, "CreateVideoUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	author, err := v.UserRepository.FindById(newCtx, cmd.AuthorId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	newPost := entity.NewVideoPost(cmd.Title, cmd.Description, cmd.VideoUUID, *author)

	err = v.VideoPostRepository.Create(newCtx, newPost)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &CreateVideoResponse{}, nil
}
