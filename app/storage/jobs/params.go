package jobs

type Params struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	RelatedId int    `json:"related_id"`
	Type      string `json:"type"`
	ProcessId string `json:"process_id"`
	Progress  int    `json:"progress"`
	Error     int    `json:"error"`
	Offset    uint64 `json:"offset"`
	Limit     uint64 `json:"limit"`
}
