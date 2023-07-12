package comment_test

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByVideoCQRS_Execute(t *testing.T) {
	type FindCommentByVideoTestCase struct {
		Description      string
		VideoId          uint
		ExpectedComments []*entity.VideoComment
		ExpectedError    error
	}

	testCases := []FindCommentByVideoTestCase{
		{
			Description: "Video has comments",
			VideoId:     1,
			ExpectedComments: []*entity.VideoComment{
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
			},
			ExpectedError: nil,
		},
		{
			Description: "Video has no comments",
			VideoId:     1,
			ExpectedComments: []*entity.VideoComment{
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
			},
			ExpectedError: nil,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d]", i)

		videoCommentRepository := mocks.NewVideoCommentRepository(t)
		videoCommentRepository.On("FindByVideoId", tc.VideoId).Return(tc.ExpectedComments, nil)

		useCase := comment.NewFindByVideoUseCase(videoCommentRepository)
		query := comment.FindByVideoQuery{VideoId: tc.VideoId}
		response, err := useCase.Execute(query)
		if tc.ExpectedError != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.EqualValues(t, tc.ExpectedComments, response.Comments)
	}
}
