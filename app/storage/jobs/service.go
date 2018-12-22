package jobs

import (
	"database/sql"
	"errors"
	"fmt"
)

type Service struct {
	repo *Repository
}

func (this *Service) AddJob(UserId int, RelatedId int, Type string) (job *Job, err error) {
	if UserId == 0 || RelatedId == 0 || Type == "" {
		return nil, errors.New("wrong arguments")
	}

	job, err = this.FindOne(Params{UserId: UserId, RelatedId: RelatedId, Type: Type})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if job != nil {
		return
	}

	err = this.Insert(&Job{
		UserId:    UserId,
		RelatedId: RelatedId,
		Type:      Type,
	})

	return
}

func (this *Service) CheckJobStatus(UserId int, RelatedId int, Type string) (status string, err error) {
	if UserId == 0 || RelatedId == 0 || Type == "" {
		return "", errors.New("wrong arguments")
	}

	job, err := this.FindOne(Params{UserId: UserId, RelatedId: RelatedId, Type: Type})
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		// check for fail job
		fJob, err := this.repo.FindAllFailedJobs(Params{UserId: UserId, RelatedId: RelatedId, Type: Type})
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if len(fJob) == 0 || err == sql.ErrNoRows {
			return "done", nil
		}

		return "failed", errors.New(fmt.Sprintf("Failed reason: %s", fJob[0].Error))
	}

	if job.ProcessId == 0 {
		return "queued", nil
	}

	return "processing", nil
}

func (this *Service) FindOne(params Params) (*Job, error) {
	params.Limit = 1
	params.Offset = 0
	_job, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_job) == 0 {
		return nil, sql.ErrNoRows
	}
	return _job[0], nil
}

func (this *Service) FindAll(params Params) ([]*Job, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Job) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Job) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Job) error {
	return this.repo.Delete(item)
}

func (this *Service) FindOneFailedJob(params Params) (*FailedJob, error) {
	params.Limit = 1
	params.Offset = 0
	_job, err := this.repo.FindAllFailedJobs(params)
	if err != nil {
		return nil, err
	}
	if len(_job) == 0 {
		return nil, sql.ErrNoRows
	}
	return _job[0], nil
}

func (this *Service) JobFailed(item *Job, errorMessage string) (err error) {
	err = this.repo.Delete(item)
	if err != nil {
		return
	}

	failedJob := item.ConvertToFailedJob(errorMessage)

	_ = this.repo.DeleteFailedJobByType(failedJob)
	return this.repo.InsertFailedJob(failedJob)
}

func (this *Service) JobCompleted(item *Job) (err error) {
	err1 := this.repo.Delete(item)
	err2 := this.repo.DeleteFailedJobByType(item.ConvertToFailedJob(""))

	if err1 != nil || err2 != nil {
		fmt.Printf("SQL error Delete Job id:'%d'.\n\tError1:\n%v\n\n\tError1:\n%v\n", item.Id, err1, err2)
		err = errors.New("JobComplete error")
	}

	return
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
