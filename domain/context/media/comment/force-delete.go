package comment

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func NewForceDeleteCommentUseCase(videoCommentRepository repository.VideoCommentRepository) *ForceDeleteCommentUseCase {
	return &ForceDeleteCommentUseCase{
		VideoCommentRepository: videoCommentRepository,
	}
}

func (uc ForceDeleteCommentUseCase) Execute(_ context.Context, cmd ForceDeleteCommentCommand) (*ForceDeleteCommentResponse, error) {
	isAdmin := false
	for _, roleId := range cmd.RoleIds {
		if roleId == role.Admin || roleId == role.SuperAdmin {
			isAdmin = true
		}
	}

	if !isAdmin {
		return nil, exception.ErrUnauthorized
	}

	deleteCount, err := uc.VideoCommentRepository.ForceDeleteById(cmd.CommentId)
	if err != nil {
		return nil, err
	}

	return &ForceDeleteCommentResponse{DeleteCount: deleteCount}, nil
}
