package video

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ ILikeVideoUseCase = (*LikeVideoUseCase)(nil)

type LikeVideoCommand struct {
	VideoId uint `json:"videoId,omitempty"`
	UserId  uint `json:"userId,omitempty"`
}

type LikeVideoResponse struct {
}

type ILikeVideoUseCase interface {
	Execute(ctx context.Context, cmd LikeVideoCommand) (*LikeVideoResponse, error)
}

type LikeVideoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	TracerProvider      *trace.TracerProvider
}

func NewLikeVideoUseCase(tracerProvider *trace.TracerProvider, videoPostRepository repository.VideoPostRepository) *LikeVideoUseCase {
	return &LikeVideoUseCase{
		VideoPostRepository: videoPostRepository,
		TracerProvider:      tracerProvider,
	}
}

func (uc LikeVideoUseCase) Execute(ctx context.Context, cmd LikeVideoCommand) (*LikeVideoResponse, error) {
	err := uc.VideoPostRepository.Like(ctx, cmd.VideoId, cmd.UserId)
	if err != nil {
		return nil, err
	}
	return &LikeVideoResponse{}, nil
}
