package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"time"
)

type CreateCommentCommand struct {
	Comment  string `json:"comment"`
	VideoId  uint   `json:"videoId,string"`
	AuthorId uint   `json:"-"`
}

type CreateCommentResponse struct {
	CommendId string `json:"commend_id"`
}

type IVideoCommentCreateUseCase interface {
	Execute(cmd CreateCommentCommand) (*CreateCommentResponse, error)
}

type CreateCommentUseCase struct {
	VideoCommentRepo    repository.VideoCommentRepository
	UserRepository      repository.UserRepository
	VideoPostRepository repository.VideoPostRepository
}

func NewCreateCommentUseCase(videoCommentRepo repository.VideoCommentRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{VideoCommentRepo: videoCommentRepo}
}

func (c CreateCommentUseCase) Execute(cmd CreateCommentCommand) (*CreateCommentResponse, error) {
	author, err := c.UserRepository.FindById(cmd.AuthorId)
	if err != nil {
		return nil, err
	}

	video, err := c.VideoPostRepository.FindById(cmd.VideoId)
	if err != nil {
		return nil, err
	}

	comment := entity.NewVideoComment(*video, *author, cmd.Comment, time.Now(), time.Now())

	err = c.VideoCommentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return &CreateCommentResponse{CommendId: comment.Id.Hex()}, nil
}
