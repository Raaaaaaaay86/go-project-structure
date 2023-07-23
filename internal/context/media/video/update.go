package video

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"gorm.io/gorm"
)

var _ validate.Validator = (*UpdateVideoInfoCommand)(nil)

type UpdateVideoInfoCommand struct {
	UpdaterId   uint
	VideoId     uint
	Title       string
	Description string
}

func (c UpdateVideoInfoCommand) Validate() error {
	validations := []struct {
		ValidatedResult bool
		Err             func() error
	}{
		{
			ValidatedResult: c.UpdaterId != 0,
			Err:             func() error { return exception.NewInvalidInputError("updaterId").ShouldNotEmpty() },
		},
		{
			ValidatedResult: c.VideoId != 0,
			Err:             func() error { return exception.NewInvalidInputError("videoId").ShouldNotEmpty() },
		},
		{
			ValidatedResult: len(c.Title) > 0,
			Err:             func() error { return exception.NewInvalidInputError("title").ShouldNotEmpty() },
		},
	}

	for _, validation := range validations {
		if !validation.ValidatedResult {
			return validation.Err()
		}
	}

	return nil
}

type UpdateVideoInfoResponse struct {
	VideoPost *entity.VideoPost
}

type IUpdateVideoInfoUseCase interface {
	Execute(ctx context.Context, cmd UpdateVideoInfoCommand) (*UpdateVideoInfoResponse, error)
}

var _ IUpdateVideoInfoUseCase = (*UpdateVideoInfoUseCase)(nil)

type UpdateVideoInfoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	DB                  *gorm.DB
	TracerProvider      tracing.ApplicationTracer
}

func NewUpdateVideoInfoUseCase(tracerProvider tracing.ApplicationTracer, videoPostRepository repository.VideoPostRepository, DB *gorm.DB) *UpdateVideoInfoUseCase {
	return &UpdateVideoInfoUseCase{
		VideoPostRepository: videoPostRepository,
		DB:                  DB,
		TracerProvider:      tracerProvider,
	}
}

func (uc UpdateVideoInfoUseCase) Execute(ctx context.Context, cmd UpdateVideoInfoCommand) (*UpdateVideoInfoResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(uc.TracerProvider, ctx, pkg, "UpdateVideoInfoUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	video, err := uc.VideoPostRepository.FindById(newCtx, cmd.VideoId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	isNotUpdatedByAuthor := video.AuthorId != cmd.UpdaterId
	if isNotUpdatedByAuthor {
		return nil, exception.ErrUnauthorized
	}

	video.Title = cmd.Title
	video.Description = cmd.Description

	err = uc.VideoPostRepository.Update(newCtx, video)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &UpdateVideoInfoResponse{VideoPost: video}, nil
}
