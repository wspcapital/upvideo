package invites

type Params struct {
	Id     int
	UserId int
	Title  string
	Code   string
	Offset uint64
	Limit  uint64
}
