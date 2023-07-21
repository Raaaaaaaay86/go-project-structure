package comment

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ IDeleteCommentUseCase = (*DeleteCommentUseCase)(nil)

type DeleteCommentCommand struct {
	CommentId  primitive.ObjectID `json:"commentId"`
	ExecutorId uint               `json:"-"`
	RoleIds    []role.RoleId      `json:"-"`
}

type DeleteCommentResponse struct {
}

type IDeleteCommentUseCase interface {
	Execute(ctx context.Context, cmd DeleteCommentCommand) (*DeleteCommentResponse, error)
}

type DeleteCommentUseCase struct {
	VideoCommentRepository repository.VideoCommentRepository
	TracerProvider         *trace.TracerProvider
}

func NewDeleteCommentUseCase(tracerProvider *trace.TracerProvider, videoCommentRepository repository.VideoCommentRepository) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{
		VideoCommentRepository: videoCommentRepository,
		TracerProvider:         tracerProvider,
	}
}

func (d DeleteCommentUseCase) Execute(ctx context.Context, cmd DeleteCommentCommand) (*DeleteCommentResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(d.TracerProvider, ctx, pkg, "DeleteCommentUseCase.Execute")
	defer span.End()

	_, err := d.VideoCommentRepository.DeleteById(newCtx, cmd.CommentId, cmd.ExecutorId)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &DeleteCommentResponse{}, nil
}
