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
	"gorm.io/gorm"
	"testing"
)

func TestCreateCommentUseCase_Execute(t *testing.T) {
	type CreateCommentUseCase struct {
		AuthorId      uint
		VideoId       uint
		Comment       string
		ExpectedError error
	}

	testCases := []CreateCommentUseCase{
		{
			AuthorId:      1,
			VideoId:       1,
			Comment:       "test comment",
			ExpectedError: nil,
		},
		{
			AuthorId:      1,
			VideoId:       1,
			Comment:       "",
			ExpectedError: exception.ErrEmptyInput,
		},
		{
			AuthorId:      1,
			VideoId:       1,
			Comment:       "test comment",
			ExpectedError: gorm.ErrRecordNotFound,
		},
		{
			AuthorId:      1,
			VideoId:       1,
			Comment:       "test comment",
			ExpectedError: gorm.ErrRecordNotFound,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d]", i)

		userRepository := mocks.NewUserRepository(t)
		videoPostRepository := mocks.NewVideoPostRepository(t)
		videoCommentRepository := mocks.NewVideoCommentRepository(t)

		switch tc.ExpectedError {
		case nil:
			videoPostRepository.On("FindById", mock.Anything, tc.VideoId).
				Return(&entity.VideoPost{Id: tc.VideoId}, tc.ExpectedError).
				Once()
			userRepository.On("FindById", mock.Anything, tc.AuthorId).
				Return(&entity.User{Id: tc.AuthorId}, tc.ExpectedError).
				Once()
			videoCommentRepository.On("Create", mock.Anything, mock.AnythingOfType("*entity.VideoComment")).
				Return(nil).
				Once()
		case gorm.ErrRecordNotFound:
			videoPostRepository.On("FindById", mock.Anything, tc.VideoId).
				Return(&entity.VideoPost{Id: tc.VideoId}, tc.ExpectedError).
				Maybe()
			userRepository.On("FindById", mock.Anything, tc.AuthorId).
				Return(&entity.User{Id: tc.AuthorId}, tc.ExpectedError).
				Maybe()
		case exception.ErrEmptyInput:
		}

		cmd := comment.CreateCommentCommand{
			AuthorId: tc.AuthorId,
			VideoId:  tc.VideoId,
			Comment:  tc.Comment,
		}
		useCase := comment.NewCreateCommentUseCase(tracing.NewEmptyTracerProvider(), videoCommentRepository, userRepository, videoPostRepository)
		response, err := useCase.Execute(context.TODO(), cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.NotEmptyf(t, response.CommendId, "response comment id should not be empty")
	}
}
