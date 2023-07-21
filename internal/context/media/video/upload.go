package video

import (
	"context"
	"github.com/google/uuid"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
	"io"
)

var _ IUploadVideoUseCase = (*UploadVideoUseCase)(nil)

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
	Execute(ctx context.Context, cmd UploadVideoCommand) (*UploadVideoResponse, error)
}

type UploadVideoUseCase struct {
	FileUploader   bucket.Uploader
	Ffmpeg         convert.IFfmpeg
	TracerProvider *trace.TracerProvider
}

func NewUploadVideoUseCase(tracerProvider *trace.TracerProvider, fileUploader bucket.Uploader, ffmpeg convert.IFfmpeg) *UploadVideoUseCase {
	return &UploadVideoUseCase{
		FileUploader:   fileUploader,
		Ffmpeg:         ffmpeg,
		TracerProvider: tracerProvider,
	}
}

func (c UploadVideoUseCase) Execute(ctx context.Context, cmd UploadVideoCommand) (*UploadVideoResponse, error) {
	_, span := tracing.ApplicationSpanFactory(c.TracerProvider, ctx, pkg, "UploadVideoUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	result, err := c.FileUploader.Upload(cmd.File, cmd.FileName)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	uid := uuid.New().String()
	err = c.Ffmpeg.Convert(result.FullPath, uid)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &UploadVideoResponse{VideoUUID: uid}, nil
}
