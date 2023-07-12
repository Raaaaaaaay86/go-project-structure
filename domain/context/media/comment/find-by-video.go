package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
)

type FindByVideoQuery struct {
	VideoId uint `form:"videoId"`
}

func (f FindByVideoQuery) Validate() error {
	if f.VideoId == 0 {
		return exception.ErrEmptyInput
	}
	return nil
}

type FindByVideoResponse struct {
	Comments []*entity.VideoComment `json:"comments"`
}

type IFindByVideoCQRS interface {
	Execute(query FindByVideoQuery) (*FindByVideoResponse, error)
}

type FindByVideoCQRS struct {
	VideoCommentRepository repository.VideoCommentRepository
}

func NewFindByVideoUseCase(videoCommentRepository repository.VideoCommentRepository) *FindByVideoCQRS {
	return &FindByVideoCQRS{VideoCommentRepository: videoCommentRepository}
}

func (f FindByVideoCQRS) Execute(query FindByVideoQuery) (*FindByVideoResponse, error) {
	err := validate.Do(query)
	if err != nil {
		return nil, err
	}

	comments, err := f.VideoCommentRepository.FindByVideoId(query.VideoId)
	if err != nil {
		return nil, err
	}
	return &FindByVideoResponse{comments}, nil
}
