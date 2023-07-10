package repository

import "github.com/raaaaaaaay86/go-project-structure/domain/entity"

type VideoCommentRepository interface {
	Create(comment *entity.VideoComment) error
	FindByVideoId(videoId uint) ([]*entity.VideoComment, error)
}
