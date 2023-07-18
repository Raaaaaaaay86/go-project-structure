package video_test

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestVideoCreateCQRS_Execute(t *testing.T) {
	type VideoCreateTestCase struct {
		Title         string
		Description   string
		VideoUUID     string
		Author        entity.User
		ExpectedError error
	}

	testCases := []VideoCreateTestCase{
		{
			Title:       "test title",
			Description: "test description",
			VideoUUID:   "test uuid",
			Author: entity.User{
				Id: 1,
			},
			ExpectedError: nil,
		},
		{
			Title:       "test title",
			Description: "test description",
			VideoUUID:   "",
			Author: entity.User{
				Id: 1,
			},
			ExpectedError: exception.ErrEmptyInput,
		},
		{
			Title:       "test title",
			Description: "",
			VideoUUID:   "test uuid",
			Author: entity.User{
				Id: 1,
			},
		},
		{
			Title:       "",
			Description: "test description",
			VideoUUID:   "test uuid",
			Author: entity.User{
				Id: 1,
			},
			ExpectedError: exception.ErrEmptyInput,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d]", i)

		cmd := video.CreateVideoCommand{
			Title:       tc.Title,
			Description: tc.Description,
			VideoUUID:   tc.VideoUUID,
			AuthorId:    tc.Author.Id,
		}

		videoRepository := mocks.NewVideoPostRepository(t)
		userRepository := mocks.NewUserRepository(t)
		switch tc.ExpectedError {
		case nil:
			userRepository.On("FindById", mock.Anything, cmd.AuthorId).Return(&tc.Author, nil).Once()

			post := entity.NewVideoPost(cmd.Title, cmd.Description, cmd.VideoUUID, tc.Author)
			videoRepository.On("Create", mock.Anything, post).Return(nil).Once()
		case exception.ErrEmptyInput:
		}

		_, err := video.NewCreateVideoUseCase(tracing.NewEmptyTracerProvider(), videoRepository, userRepository).Execute(context.TODO(), cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}
	}
}
