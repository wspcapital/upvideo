package shortlinks

import (
	"time"
)

type Params struct {
	Id          int        `json:"id"`
	UserId      int        `json:"user_id"`
	UniqId      string     `json:"uniq_id"`
	Url         string     `json:"url"`
	Counter     int64      `json:"counter"`
	Disabled    bool       `json:"disabled"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Offset     uint64      `json:"offset"`
	Limit      uint64      `json:"limit"`
}