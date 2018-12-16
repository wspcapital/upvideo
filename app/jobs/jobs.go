package jobs

import "time"

type Job struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	RelatedId int        `json:"related_id"`
	Type      string     `json:"type"`
	ProcessId int        `json:"process_id"`
	Progress  int        `json:"progress"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type FailedJob struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	RelatedId int        `json:"related_id"`
	Type      string     `json:"type"`
	ProcessId int        `json:"process_id"`
	Error     string     `json:"error"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (this *Job) ConvertToFailedJob(errorMessage string) *FailedJob {
	return &FailedJob{
		UserId:    this.UserId,
		RelatedId: this.RelatedId,
		Type:      this.Type,
		ProcessId: this.ProcessId,
		Error:     errorMessage,
	}
}
