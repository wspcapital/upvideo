package titles

import (
	"fmt"
	"strings"
)

type Title struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	CampaignId int    `json:"campaign_id"`
	Title      string `json:"title"`
	Tags       string `json:"tags"`
	File       string `json:"file"`
	TmpFile    string `json:"tmp_file"`
	YoutubeId  string `json:"youtube_id"`
	Posted     bool   `json:"posted"`
	Converted  bool   `json:"converted"`
	Pending    bool   `json:"pending"`
	FrameRate  int    `json:"frame_rate"`
	Resolution int    `json:"resolution"`
	IpAddress  string `json:"ip_address"`
}

func (this *Title) GetPreparedFilename() string {
	if this.Title == "" {
		fmt.Println("Title GetPreparedFilename of empty Title")
		return "empty-title.mp4"
	}
	return strings.Replace(this.Title, " ", "-", -1) + ".mp4"
}
