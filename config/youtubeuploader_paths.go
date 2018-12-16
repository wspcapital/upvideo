package config

import (
	"fmt"
	"os"
	"path"
)

type YoutubeUploaderDirs struct {
	SecretsDir      string
	TokensDir       string
	TempVideosDir   string
	MetasDir        string
	TitlesVideosDir string
}

func (this *YoutubeUploaderDirs) resolvePaths(youtubeUploaderPath string) (err error) {
	this.SecretsDir = path.Join(youtubeUploaderPath, "/secrets")
	if err = os.MkdirAll(this.SecretsDir, os.ModePerm); err != nil {
		fmt.Println("\n Create SecretsDir Error: ", err.Error())
		return
	}

	this.TokensDir = path.Join(youtubeUploaderPath, "/tokens")
	if err = os.MkdirAll(this.TokensDir, os.ModePerm); err != nil {
		fmt.Println("\n Create TokensDir Error: ", err.Error())
		return
	}

	this.TempVideosDir = path.Join(youtubeUploaderPath, "/temp_videos")
	if err = os.MkdirAll(this.TempVideosDir, os.ModePerm); err != nil {
		fmt.Println("\n Create TempVideosDir Error: ", err.Error())
		return
	}

	this.MetasDir = path.Join(youtubeUploaderPath, "/metas")
	if err = os.MkdirAll(this.MetasDir, os.ModePerm); err != nil {
		fmt.Println("\n Create MetasDir Error: ", err.Error())
		return
	}

	this.TitlesVideosDir = path.Join(youtubeUploaderPath, "/titles_videos")
	if err = os.MkdirAll(this.TitlesVideosDir, os.ModePerm); err != nil {
		fmt.Println("\n Create TitlesVideosDir Error: ", err.Error())
		return
	}

	return
}
