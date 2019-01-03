package topusers

import (
	"time"
)

type TopUser struct {
	Id               int        `json:"id"`
	Email            string     `json:"author_username"`
	NewVideosCount   int        `json:"new_videos"`
	NewChannelsCount int        `json:"new_channels"`
	LastActivity     *time.Time `json:"last_activity"`
}
