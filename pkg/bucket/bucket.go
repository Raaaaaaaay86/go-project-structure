package bucket

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type UploadResult struct {
	FileName      string
	FileExtension string
	FullPath      string
}

type Uploader interface {
	Upload(file io.Reader, fileName string) (*UploadResult, error)
}

type LocalUploader struct {
	BucketPath string
}

func NewLocalUploader(destination string) *LocalUploader {
	return &LocalUploader{
		destination,
	}
}

func (l LocalUploader) Upload(file io.Reader, fileName string) (*UploadResult, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// "/Users/linjiafu/project/go-video-demo/bucket/%s"
	filePath := fmt.Sprintf("%s/%s", l.BucketPath, fileName)
	targetFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	_, err = targetFile.Write(data)
	if err != nil {
		return nil, err
	}

	s := strings.Split(fileName, ".")
	name := strings.Join(s[:len(s)-1], "")
	extension := s[len(s)-1]

	return &UploadResult{
		FileName:      name,
		FileExtension: extension,
		FullPath:      filePath,
	}, nil
}
