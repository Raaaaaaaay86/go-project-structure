package video_test

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestUpdateVideoInfoUseCase_Execute(t *testing.T) {
	type UpdateVideoInfoTestCase struct {
		TestDescription   string
		UpdateTitle       string
		UpdateDescription string
		UpdaterId         uint
		VideoId           uint
		VideoAuthorId     uint
		ExpectedError     error
	}

	testCases := []UpdateVideoInfoTestCase{
		{
			TestDescription:   "Missing UpdaterId",
			UpdateTitle:       "New title",
			UpdateDescription: "New description",
			UpdaterId:         0,
			VideoId:           2,
			ExpectedError:     exception.ErrEmptyInput,
		},
		{
			TestDescription:   "Missing VideoId",
			UpdateTitle:       "New title",
			UpdateDescription: "New description",
			UpdaterId:         1,
			VideoAuthorId:     1,
			VideoId:           0,
			ExpectedError:     exception.ErrEmptyInput,
		},
		{
			TestDescription:   "Missing Title",
			UpdateTitle:       "",
			UpdateDescription: "New description",
			UpdaterId:         1,
			VideoAuthorId:     1,
			VideoId:           2,
			ExpectedError:     exception.ErrEmptyInput,
		},
		{
			TestDescription:   "Update by non-author",
			UpdateTitle:       "New title",
			UpdateDescription: "New description",
			UpdaterId:         2,
			VideoAuthorId:     1,
			VideoId:           2,
			ExpectedError:     exception.ErrUnauthorized,
		},
		{
			TestDescription:   "Update by author",
			UpdateTitle:       "New title",
			UpdateDescription: "New description",
			UpdaterId:         1,
			VideoAuthorId:     1,
			VideoId:           2,
			ExpectedError:     nil,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d]: %s", i, tc.TestDescription)

		db := &gorm.DB{}
		videoPostRepository := mocks.NewVideoPostRepository(t)
		videoPostRepository.On("StartTx").Return(db).Maybe()
		videoPostRepository.On("CommitTx", db).Return(db).Maybe()
		videoPostRepository.On("WithTx", db).Return(videoPostRepository).Maybe()
		videoPostRepository.On("ForUpdate").Return(videoPostRepository).Maybe()

		oldVideo := &entity.VideoPost{Id: tc.VideoId, AuthorId: tc.VideoAuthorId}

		updatedVideo := &entity.VideoPost{
			Id:          oldVideo.Id,
			AuthorId:    oldVideo.AuthorId,
			Title:       tc.UpdateTitle,
			Description: tc.UpdateDescription,
		}

		switch tc.ExpectedError {
		case nil:
			videoPostRepository.On("FindById", tc.VideoId).Return(oldVideo, nil).Once()
			videoPostRepository.On("Update", updatedVideo).Return(nil).Once()
		case exception.ErrUnauthorized:
			videoPostRepository.On("FindById", tc.VideoId).Return(oldVideo, nil).Once()
		case exception.ErrEmptyInput:
		}

		cmd := video.UpdateVideoInfoCommand{
			UpdaterId:   tc.UpdaterId,
			VideoId:     tc.VideoId,
			Title:       tc.UpdateTitle,
			Description: tc.UpdateDescription,
		}
		response, err := video.NewUpdateVideoInfoUseCase(videoPostRepository, db).Execute(cmd)
		if err != nil || tc.ExpectedError != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.Equal(t, tc.UpdateTitle, response.VideoPost.Title)
		assert.Equal(t, tc.UpdateDescription, response.VideoPost.Description)
	}
}
