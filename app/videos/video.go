package videos

type Video struct {
	Id          int
	UserId      int
	AccountId   int
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
}
