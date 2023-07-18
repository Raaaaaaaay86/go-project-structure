package video

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

type UpdateVideoInfoCommand struct {
	UpdaterId   uint
	VideoId     uint
	Title       string
	Description string
}

func (c UpdateVideoInfoCommand) Validate() error {
	checkList := []bool{
		c.UpdaterId != 0,
		c.VideoId != 0,
		c.Title != "",
	}

	for _, ok := range checkList {
		if !ok {
			return exception.ErrEmptyInput
		}
	}

	return nil
}

type UpdateVideoInfoResponse struct {
	VideoPost *entity.VideoPost
}

type IUpdateVideoInfoUseCase interface {
	Execute(cmd UpdateVideoInfoCommand) (*UpdateVideoInfoResponse, error)
}

type UpdateVideoInfoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	DB                  *gorm.DB
	TracerProvider      *trace.TracerProvider
}

func NewUpdateVideoInfoUseCase(tracerProvider *trace.TracerProvider, videoPostRepository repository.VideoPostRepository, DB *gorm.DB) *UpdateVideoInfoUseCase {
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

	tx := uc.VideoPostRepository.StartTx()

	video, err := uc.VideoPostRepository.WithTx(tx).FindById(newCtx, cmd.VideoId)
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

	err = uc.VideoPostRepository.WithTx(tx).ForUpdate().Update(newCtx, video)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	uc.VideoPostRepository.CommitTx(tx)

	return &UpdateVideoInfoResponse{VideoPost: video}, nil
}
