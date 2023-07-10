package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteCommentCommand struct {
	CommentId  primitive.ObjectID `json:"commentId"`
	ExecutorId uint               `json:"executorId"`
	RoleId     role.RoleId        `json:"roleId"`
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
	comment, err := d.VideoCommentRepository.FindById(cmd.CommentId)
	if err != nil {
		return nil, err
	}

	isNotAuthor := comment.AuthorId != cmd.ExecutorId
	// There would be a better way to do RBAC checking
	isNotAdmin := cmd.RoleId != role.Admin || cmd.RoleId != role.SuperAdmin
	if isNotAuthor && isNotAdmin {
		return nil, exception.ErrUnauthorized
	}

	err = d.VideoCommentRepository.DeleteById(cmd.CommentId)
	if err != nil {
		return nil, err
	}

	return &DeleteCommentResponse{}, nil
}
