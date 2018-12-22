package titles

type Params struct {
	Id         int
	UserId     int
	CampaignId int
	Title      string
	Tags       string
	File       string
	TmpFile    string
	YoutubeId  string
	Posted     bool
	Converted  bool
	Pending    bool
	FrameRate  int
	Resolution int
	IpAddress  string
	Offset     uint64
	Limit      uint64
}
