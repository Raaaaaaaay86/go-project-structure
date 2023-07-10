package video_test

import (
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVideoUploadCQRS_Execute(t *testing.T) {
	type VideoUploadTestCase struct {
		File        io.Reader
		FileName    string
		ExpectedErr error
	}

	createTempFile := func() *os.File {
		t.Helper()
		file, err := os.CreateTemp("", "test*.mp4")
		if err != nil {
			t.Fatal(err)
		}
		return file
	}

	getExtension := func(fileName string) string {
		t.Helper()
		s := strings.Split(fileName, ".")
		return s[len(s)-1]
	}

	getName := func(fileName string) string {
		t.Helper()
		s := strings.Split(fileName, ".")
		return strings.Join(s[:len(s)-1], "")
	}

	testCases := []VideoUploadTestCase{
		{
			File:        nil,
			FileName:    "hello.mp4",
			ExpectedErr: exception.ErrEmptyFile,
		},
		{
			File:        createTempFile(),
			FileName:    "",
			ExpectedErr: exception.ErrEmptyInput,
		},
		{
			File:        createTempFile(),
			FileName:    "world.mp4",
			ExpectedErr: nil,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d]", i)

		cmd := video.UploadVideoCommand{
			File:     tc.File,
			FileName: tc.FileName,
		}

		uploader := mocks.NewUploader(t)
		ffmpeg := mocks.NewIFfmpeg(t)
		switch tc.ExpectedErr {
		case nil:
			expectedResult := &bucket.UploadResult{
				FullPath:      fmt.Sprintf("path/to/%s", cmd.FileName),
				FileName:      getName(cmd.FileName),
				FileExtension: getExtension(cmd.FileName),
			}

			uploader.On("Upload", cmd.File, cmd.FileName).Return(expectedResult, nil).Once()
			ffmpeg.On("Convert", mock.Anything, mock.Anything).Return(nil).Once()
		case exception.ErrEmptyInput, exception.ErrEmptyFile:
		}

		response, err := video.NewUploadVideoUseCase(uploader, ffmpeg).Execute(cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedErr, err)
			continue
		}

		assert.NotEmpty(t, response.VideoUUID)
	}
}
