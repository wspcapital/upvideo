package topusers

import (
	"time"
)

type TopUser struct {
	Id                int        `json:"id"`
	Email             string     `json:"author_username"`
	NewAccountsCount  int        `json:"new_accounts"`
	NewCampaignsCount int        `json:"new_campaigns"`
	LastActivity      *time.Time `json:"last_activity"`
}
