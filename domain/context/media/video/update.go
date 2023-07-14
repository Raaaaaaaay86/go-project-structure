package video

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"gorm.io/gorm"
)

type UpdateVideoInfoCommand struct {
	UpdaterId   uint
	VideoId     uint
	Title       string
	Description string
}

func (c UpdateVideoInfoCommand) Validate() error {
	checkList := []bool{
		c.UpdaterId != 0,
		c.VideoId != 0,
		c.Title != "",
	}

	for _, ok := range checkList {
		if !ok {
			return exception.ErrEmptyInput
		}
	}

	return nil
}

type UpdateVideoInfoResponse struct {
	VideoPost *entity.VideoPost
}

type IUpdateVideoInfoUseCase interface {
	Execute(cmd UpdateVideoInfoCommand) (*UpdateVideoInfoResponse, error)
}

type UpdateVideoInfoUseCase struct {
	VideoPostRepository repository.VideoPostRepository
	DB                  *gorm.DB
}

func NewUpdateVideoInfoUseCase(videoPostRepository repository.VideoPostRepository, DB *gorm.DB) *UpdateVideoInfoUseCase {
	return &UpdateVideoInfoUseCase{
		VideoPostRepository: videoPostRepository,
		DB:                  DB,
	}
}

func (uc UpdateVideoInfoUseCase) Execute(cmd UpdateVideoInfoCommand) (*UpdateVideoInfoResponse, error) {
	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	tx := uc.VideoPostRepository.StartTx()

	video, err := uc.VideoPostRepository.WithTx(tx).FindById(cmd.VideoId)
	if err != nil {
		return nil, err
	}

	isNotUpdatedByAuthor := video.AuthorId != cmd.UpdaterId
	if isNotUpdatedByAuthor {
		return nil, exception.ErrUnauthorized
	}

	video.Title = cmd.Title
	video.Description = cmd.Description

	err = uc.VideoPostRepository.WithTx(tx).ForUpdate().Update(video)
	if err != nil {
		return nil, err
	}

	uc.VideoPostRepository.CommitTx(tx)

	return &UpdateVideoInfoResponse{VideoPost: video}, nil
}
