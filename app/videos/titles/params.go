package titles

type Params struct {
	Id         int
	UserId     int
	VideoId    int
	Title      string
	Tags       string
	File       string
	Posted     bool
	Converted  bool
	Pending    bool
	FrameRate  int
	Resolution int
	IpAddress  string
	Offset     uint64
	Limit      uint64
}
