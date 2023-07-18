package comment

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ForceDeleteCommentCommand struct {
	CommentId  primitive.ObjectID `json:"commentId"`
	ExecutorId uint               `json:"-"`
	RoleIds    []role.RoleId      `json:"-"`
}

type ForceDeleteCommentResponse struct {
	DeleteCount int `json:"deleteCount"`
}

type IForceDeleteCommentUseCase interface {
	Execute(ctx context.Context, cmd ForceDeleteCommentCommand) (*ForceDeleteCommentResponse, error)
}

type ForceDeleteCommentUseCase struct {
	VideoCommentRepository repository.VideoCommentRepository
	TracerProvider         *trace.TracerProvider
}

func NewForceDeleteCommentUseCase(tracerProvider *trace.TracerProvider, videoCommentRepository repository.VideoCommentRepository) *ForceDeleteCommentUseCase {
	return &ForceDeleteCommentUseCase{
		VideoCommentRepository: videoCommentRepository,
		TracerProvider:         tracerProvider,
	}
}

func (uc ForceDeleteCommentUseCase) Execute(ctx context.Context, cmd ForceDeleteCommentCommand) (*ForceDeleteCommentResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(uc.TracerProvider, ctx, pkg, "ForceDeleteCommentUseCase.Execute")
	defer span.End()

	isAdmin := false
	for _, roleId := range cmd.RoleIds {
		if roleId == role.Admin || roleId == role.SuperAdmin {
			isAdmin = true
		}
	}

	if !isAdmin {
		return nil, exception.ErrUnauthorized
	}

	deleteCount, err := uc.VideoCommentRepository.ForceDeleteById(newCtx, cmd.CommentId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &ForceDeleteCommentResponse{DeleteCount: deleteCount}, nil
}
