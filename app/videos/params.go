package videos

type Params struct {
	Id          int
	UserId      int
	Title       string
	Description string
	Tags        string
	Category    string
	Language    string
	File        string
	Playlist    string
	IpAddress   string
	Offset      uint64
	Limit       uint64
}
