package convert

import (
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	"os/exec"
)

type IFfmpeg interface {
	Convert(targetPath string, outputFileName string) error
}

type Ffmpeg struct {
	BucketPath string
}

func NewFfmpeg(config *configs.YamlConfig) *Ffmpeg {
	return &Ffmpeg{
		BucketPath: config.BucketPath.Converted,
	}
}

func (f Ffmpeg) Convert(targetPath string, outputFileName string) error {
	// ffmpeg -i $TargetPath -c:v libx264 -c:a aac -f hls -hls_list_size 0 $OutputPath
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		targetPath,
		"-c:v",
		"libx264",
		"-c:a",
		"aac",
		"-f",
		"hls",
		"-hls_list_size",
		"0",
		fmt.Sprintf("%s/%s.m3u8", f.BucketPath, outputFileName),
	)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
