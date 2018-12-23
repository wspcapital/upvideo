package invites

type Params struct {
	Id     int        `json:"id"`
	UserId int        `json:"user_id"`
	Title  string     `json:"title"`
	Code   string     `json:"code"`
	Offset uint64     `json:"offset"`
	Limit  uint64     `json:"limit"`
}
