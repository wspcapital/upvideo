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
	TmpFile     string
	Playlist    string
	TitleGen    bool
	IpAddress   string
	Offset      uint64
	Limit       uint64
}
