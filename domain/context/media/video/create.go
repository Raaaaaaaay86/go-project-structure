package video

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
)

type CreateVideoCommand struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	VideoUUID   string `json:"uuid,omitempty"`
	AuthorId    uint   `json:"authorId,omitempty"`
}

func (c CreateVideoCommand) Validate() error {
	if c.AuthorId == 0 {
		return exception.ErrEmptyInput
	}

	fields := []string{c.Title, c.VideoUUID}
	for _, field := range fields {
		if len(field) == 0 {
			return exception.ErrEmptyInput
		}
	}
	return nil
}

type CreateVideoResponse struct {
}

type IVideoCreateUseCase interface {
	Execute(cmd CreateVideoCommand) (*CreateVideoResponse, error)
}

type VideoCreateUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	UserRepository      repository.UserRepository
}

func NewCreateVideoUseCase(videoPostRepository repository.VideoPostRepository, userRepository repository.UserRepository) *VideoCreateUseCase {
	return &VideoCreateUseCase{
		VideoPostRepository: videoPostRepository,
		UserRepository:      userRepository,
	}
}

func (v VideoCreateUseCase) Execute(cmd CreateVideoCommand) (*CreateVideoResponse, error) {
	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	author, err := v.UserRepository.FindById(cmd.AuthorId)
	if err != nil {
		return nil, err
	}

	newPost := entity.NewVideoPost(cmd.Title, cmd.Description, cmd.VideoUUID, *author)

	err = v.VideoPostRepository.Create(newPost)
	if err != nil {
		return nil, err
	}

	return &CreateVideoResponse{}, nil
}
