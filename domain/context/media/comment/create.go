package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"time"
)

type CreateCommentCommand struct {
	Comment  string `json:"comment"`
	VideoId  uint   `json:"videoId,string"`
	AuthorId uint   `json:"-"`
}

func (c CreateCommentCommand) Validate() error {
	if c.Comment == "" {
		return exception.ErrEmptyInput
	}
	return nil
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

func NewCreateCommentUseCase(videoCommentRepo repository.VideoCommentRepository, userRepository repository.UserRepository, postRepository repository.VideoPostRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		VideoCommentRepo:    videoCommentRepo,
		UserRepository:      userRepository,
		VideoPostRepository: postRepository,
	}
}

func (c CreateCommentUseCase) Execute(cmd CreateCommentCommand) (*CreateCommentResponse, error) {
	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

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
