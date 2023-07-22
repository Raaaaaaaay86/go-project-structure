package comment_test

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestCreateCommentUseCase_Execute(t *testing.T) {
	type CreateCommentUseCase struct {
		TestDescription string
		AuthorId        uint
		VideoId         uint
		Comment         string
		ExpectedError   error
	}

	testCases := []CreateCommentUseCase{
		{
			TestDescription: "Success create comment",
			AuthorId:        1,
			VideoId:         1,
			Comment:         "test comment",
			ExpectedError:   nil,
		},
		{
			TestDescription: "Failed by empty comment",
			AuthorId:        1,
			VideoId:         1,
			Comment:         "",
			ExpectedError:   exception.NewInvalidInputError("comment").ShouldNotEmpty(),
		},
		{
			TestDescription: "Failed by non-exist author (user id)",
			AuthorId:        999,
			VideoId:         1,
			Comment:         "test comment",
			ExpectedError:   gorm.ErrRecordNotFound,
		},
		{
			TestDescription: "Failed by non-exist video post (video id)",
			AuthorId:        1,
			VideoId:         999,
			Comment:         "test comment",
			ExpectedError:   gorm.ErrRecordNotFound,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d] - %s", i, tc.TestDescription)

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
		case exception.InvalidInputError{}:
		}

		cmd := comment.CreateCommentCommand{
			AuthorId: tc.AuthorId,
			VideoId:  tc.VideoId,
			Comment:  tc.Comment,
		}
		useCase := comment.NewCreateCommentUseCase(tracing.NewEmptyTracerProvider(), videoCommentRepository, userRepository, videoPostRepository)
		response, err := useCase.Execute(context.TODO(), cmd)
		if err != nil {
			assert.ErrorAs(t, tc.ExpectedError, &err)
			assert.Equal(t, tc.ExpectedError.Error(), err.Error())
			continue
		}

		assert.NotEmptyf(t, response.CommendId, "response comment id should not be empty")
	}
}
