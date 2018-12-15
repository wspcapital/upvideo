package config

import (
	"encoding/json"
	"fmt"
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
	Session             SessionConfig
	WebServer           WebServerParams
	AWS                 AWSParams `json:"AWS"`
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
