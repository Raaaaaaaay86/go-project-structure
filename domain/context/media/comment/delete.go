package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteCommentCommand struct {
	CommentId  primitive.ObjectID `json:"commentId"`
	ExecutorId uint               `json:"-"`
	RoleIds    []role.RoleId      `json:"-"`
}

type DeleteCommentResponse struct {
}

type IDeleteCommentUseCase interface {
	Execute(cmd DeleteCommentCommand) (*DeleteCommentResponse, error)
}

type DeleteCommentUseCase struct {
	VideoCommentRepository repository.VideoCommentRepository
}

func NewDeleteCommentUseCase(videoCommentRepository repository.VideoCommentRepository) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{
		VideoCommentRepository: videoCommentRepository,
	}
}

func (d DeleteCommentUseCase) Execute(cmd DeleteCommentCommand) (*DeleteCommentResponse, error) {
	_, err := d.VideoCommentRepository.DeleteById(cmd.CommentId, cmd.ExecutorId)
	if err != nil {
		return nil, err
	}

	return &DeleteCommentResponse{}, nil
}
