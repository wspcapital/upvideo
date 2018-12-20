package titles

import (
	"fmt"
	"strings"
	"time"
)

var (
	initialResolution = 1278
	minFrameRate      = 20
	maxFrameRate      = 29
)

type Title struct {
	Id         int        `json:"id"`
	UserId     int        `json:"user_id"`
	CampaignId int        `json:"campaign_id"`
	Title      string     `json:"title"`
	Tags       string     `json:"tags"`
	File       string     `json:"file"`
	TmpFile    string     `json:"tmp_file"`
	YoutubeId  string     `json:"youtube_id"`
	YoutubeUrl string     `json:"youtube_url"`
	Posted     bool       `json:"posted"`
	Converted  bool       `json:"converted"`
	Pending    bool       `json:"pending"`
	FrameRate  int        `json:"frame_rate"`
	Resolution int        `json:"resolution"`
	IpAddress  string     `json:"ip_address"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (this *Title) GetPreparedFilename() string {
	if this.Title == "" {
		fmt.Println("Title GetPreparedFilename of empty Title")
		return "empty-title.mp4"
	}
	return strings.Replace(this.Title, " ", "-", -1) + ".mp4"
}

func (this *Title) SetNextFrameRate() *Title {
	if this.Resolution == 0 || this.FrameRate == 0 {
		this.Resolution = initialResolution
		this.FrameRate = minFrameRate
	} else {
		this.FrameRate++
		if this.FrameRate > maxFrameRate {
			this.FrameRate = minFrameRate
			this.Resolution += 16
		}
	}

	return this
}
