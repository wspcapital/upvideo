package invites

type Invite struct {
	Id     int       `json:"id"`
	UserId int       `json:"user_id"`
	Title  string    `json:"title" binding:"exists,alphanum,min=4,max=255"`
	Code   string    `json:"code"`
}
