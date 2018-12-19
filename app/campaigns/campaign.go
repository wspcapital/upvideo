package campaigns

import "time"

type Campaign struct {
	Id              int        `json:"id"`
	UserId          int        `json:"user_id"`
	AccountId       int        `json:"account_id"`
	VideoId         int        `json:"video_id"`
	Title           string     `json:"title"`
	TotalTitles     int        `json:"total_titles"`
	CompleteTitles  int        `json:"complete_titles"`
	TitlesGenerated bool       `json:"titles_generated"`
	IpAddress       string     `json:"ip_address"`
	DateStart_at    *time.Time `json:"date_start_at"`
	DateComplete_at *time.Time `json:"date_complete_at"`
	Created_at      *time.Time `json:"created_at"`
	Updated_at      *time.Time `json:"updated_at"`
}
