package shortlinks

import (
	"time"
)

type Shortlink struct {
	Id          int        `json:"id"`
	UserId      int        `json:"user_id"`
	UniqId      string     `json:"uniq_id"`
	Url         string     `json:"url"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
