package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type SessionConfig struct {
	Storage    string
	TTLMinutes int32
}

type Config struct {
	Dsn                 string
	CheckInterval       int
	YoutubeUploaderPath string `json:"youtubeuploader_path"`
	YoutubeUploaderCmd  string `json:"youtubeuploader_cmd"`
	TestVideoPath       string `json:"testvideo"`
	TestVideoMetaPath   string `json:"testvideo_meta"`
	YoutubeUploaderDirs YoutubeUploaderDirs
	Session             SessionConfig
	WebServer           WebServerParams
	AWS                 AWSParams         `json:"AWS"`
	Keywordtool         KeywordtoolParams `json:"keywordtool.io"`
}

func ReadConfig(filename string) Config {
	file, _ := os.Open(filename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	resolveRelativePaths(&configuration)

	ytPaths := &YoutubeUploaderDirs{}
	err = ytPaths.resolvePaths(configuration.YoutubeUploaderPath)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}

func resolveRelativePaths(conf *Config) {
	p, err := filepath.Abs(conf.YoutubeUploaderPath)
	if err != nil {
		fmt.Sprintf("Can not resolve path: %s", conf.YoutubeUploaderPath)
	}

	conf.YoutubeUploaderPath = p

	p, err = filepath.Abs(conf.YoutubeUploaderCmd)
	if err != nil {
		fmt.Sprintf("Can not resolve path: %s", conf.YoutubeUploaderCmd)
	}

	conf.YoutubeUploaderCmd = p

	p, err = filepath.Abs(conf.TestVideoPath)
	if err != nil {
		fmt.Sprintf("Can not resolve path: %s", conf.TestVideoPath)
	}

	conf.TestVideoPath = p

	p, err = filepath.Abs(conf.TestVideoMetaPath)
	if err != nil {
		fmt.Sprintf("Can not resolve path: %s", conf.TestVideoMetaPath)
	}

	conf.TestVideoMetaPath = p
}
