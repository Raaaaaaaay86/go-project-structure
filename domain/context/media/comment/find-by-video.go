package comment

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
)

type FindByVideoQuery struct {
	VideoId uint `form:"videoId"`
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
	comments, err := f.VideoCommentRepository.FindByVideoId(query.VideoId)
	if err != nil {
		return nil, err
	}
	return &FindByVideoResponse{comments}, nil
}
