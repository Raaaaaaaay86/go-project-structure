package comment_test

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindByVideoCQRS_Execute(t *testing.T) {
	type FindCommentByVideoTestCase struct {
		TestDescription  string
		VideoId          uint
		ExpectedComments []*entity.VideoComment
		ExpectedError    error
	}

	testCases := []FindCommentByVideoTestCase{
		{
			TestDescription: "Video has comments",
			VideoId:         1,
			ExpectedComments: []*entity.VideoComment{
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
			},
			ExpectedError: nil,
		},
		{
			TestDescription: "Video has no comments",
			VideoId:         1,
			ExpectedComments: []*entity.VideoComment{
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
				{Id: primitive.NewObjectID()},
			},
			ExpectedError: nil,
		},
		{
			TestDescription: "Request by empty VideoId",
			VideoId:         0,
			ExpectedError:   exception.ErrEmptyInput,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d] - %s", i, tc.TestDescription)

		videoCommentRepository := mocks.NewVideoCommentRepository(t)
		switch tc.ExpectedError {
		case nil:
			videoCommentRepository.On("FindByVideoId", mock.Anything, tc.VideoId).Return(tc.ExpectedComments, nil)
		case exception.ErrEmptyInput:
		}

		useCase := comment.NewFindByVideoUseCase(tracing.NewEmptyTracerProvider(), videoCommentRepository)
		query := comment.FindByVideoQuery{VideoId: tc.VideoId}
		response, err := useCase.Execute(context.TODO(), query)
		if tc.ExpectedError != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.EqualValues(t, tc.ExpectedComments, response.Comments)
	}
}
