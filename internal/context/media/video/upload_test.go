package video_test

import (
	"context"
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVideoUploadUseCase_Execute(t *testing.T) {
	type VideoUploadTestCase struct {
		TestDescription string
		File            io.Reader
		FileName        string
		UploaderId      uint
		ExpectedErr     error
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
			TestDescription: "Failed by empty file",
			File:            nil,
			FileName:        "hello.mp4",
			UploaderId:      1,
			ExpectedErr:     exception.NewInvalidInputError("file").ShouldNotEmpty(),
		},
		{
			TestDescription: "Failed by empty file name",
			File:            createTempFile(),
			FileName:        "",
			UploaderId:      1,
			ExpectedErr:     exception.NewInvalidInputError("fileName").ShouldNotEmpty(),
		},
		{
			TestDescription: "Success to upload video",
			File:            createTempFile(),
			FileName:        "world.mp4",
			UploaderId:      1,
			ExpectedErr:     nil,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d] - %s", i, tc.TestDescription)

		cmd := video.UploadVideoCommand{
			File:       tc.File,
			FileName:   tc.FileName,
			UploaderId: tc.UploaderId,
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
		case exception.InvalidInputError{}:
		}

		ctx := context.Background()
		response, err := video.NewUploadVideoUseCase(tracing.NewEmptyTracerProvider(), uploader, ffmpeg).Execute(ctx, cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedErr, err)
			continue
		}

		assert.NotEmpty(t, response.VideoUUID)
	}
}
