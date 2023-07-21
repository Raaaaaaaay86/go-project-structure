package video

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ IUnLikeVideoUseCase = (*UnLikeVideoUseCase)(nil)

type UnLikeVideoCommand struct {
	VideoId uint
	UserId  uint
}

type UnLikeVideoResponse struct {
}

type IUnLikeVideoUseCase interface {
	Execute(ctx context.Context, cmd UnLikeVideoCommand) (*UnLikeVideoResponse, error)
}

type UnLikeVideoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	TracerProvider      *trace.TracerProvider
}

func NewUnLikeVideoUseCase(tracerProvider *trace.TracerProvider, videoPostRepository repository.VideoPostRepository) *UnLikeVideoUseCase {
	return &UnLikeVideoUseCase{
		VideoPostRepository: videoPostRepository,
		TracerProvider:      tracerProvider,
	}
}

func (uc UnLikeVideoUseCase) Execute(ctx context.Context, cmd UnLikeVideoCommand) (*UnLikeVideoResponse, error) {
	err := uc.VideoPostRepository.UnLike(ctx, cmd.VideoId, cmd.UserId)
	if err != nil {
		return nil, err
	}
	return &UnLikeVideoResponse{}, nil
}
