package video

import (
	"github.com/google/uuid"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"io"
)

type UploadVideoCommand struct {
	File       io.Reader `json:"-"`
	FileName   string    `json:"-"`
	UploaderId uint      `json:"uploaderId,omitempty"`
}

func (c UploadVideoCommand) Validate() error {
	if c.File == nil {
		return exception.ErrEmptyFile
	}
	if len(c.FileName) <= 0 {
		return exception.ErrEmptyInput
	}
	return nil
}

type UploadVideoResponse struct {
	VideoUUID string `json:"uuid,omitempty"`
}

type IUploadVideoUseCase interface {
	Execute(cmd UploadVideoCommand) (*UploadVideoResponse, error)
}

type UploadVideoUseCase struct {
	FileUploader bucket.Uploader
	Ffmpeg       convert.IFfmpeg
}

func NewUploadVideoUseCase(fileUploader bucket.Uploader, ffmpeg convert.IFfmpeg) *UploadVideoUseCase {
	return &UploadVideoUseCase{
		FileUploader: fileUploader,
		Ffmpeg:       ffmpeg,
	}
}

func (c UploadVideoUseCase) Execute(cmd UploadVideoCommand) (*UploadVideoResponse, error) {
	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	result, err := c.FileUploader.Upload(cmd.File, cmd.FileName)
	if err != nil {
		return nil, err
	}

	uid := uuid.New().String()
	err = c.Ffmpeg.Convert(result.FullPath, uid)
	if err != nil {
		return nil, err
	}

	return &UploadVideoResponse{VideoUUID: uid}, nil
}
