package titles

import (
	"fmt"
	"strings"
)

type Title struct {
	Id         int
	UserId     int
	VideoId    int
	Title      string
	Tags       string
	File       string
	TmpFile    string
	YoutubeId  string
	Posted     bool
	Converted  bool
	Pending    bool
	FrameRate  int
	Resolution int
	IpAddress  string
}

func (this *Title) GetPreparedFilename() string {
	if this.Title == "" {
		fmt.Println("Title GetPreparedFilename of empty Title")
		return "empty-title.mp4"
	}
	return strings.Replace(this.Title, " ", "-", -1) + ".mp4"
}
