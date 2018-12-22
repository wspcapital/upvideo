package videos

type Video struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	AccountId   int    `json:"account_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	File        string `json:"file"`
	TmpFile     string `json:"tmp_file"`
	Playlist    string `json:"playlist"`
	TitleGen    bool   `json:"title_gen"`
	IpAddress   string `json:"ip_address"`
}
