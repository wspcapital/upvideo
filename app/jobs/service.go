package jobs

import (
	"database/sql"
	"errors"
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
		if err == sql.ErrNoRows {
			return "done", nil
		}

		return "failed", errors.New("Failed reason: " + fJob[0].Error)
	}

	if job.ProcessId == "" {
		return "queued", nil
	}

	return "processing", nil
}

func (this *Service) FindOne(params Params) (*Job, error) {
	params.Limit = 1
	params.Offset = 0
	_account, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_account) == 0 {
		return nil, errors.New("found no matching accounts")
	}
	return _account[0], nil
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

func (this *Service) JobFailed(item *Job, errorMessage string) (err error) {
	err = this.repo.Delete(item)
	if err != nil {
		return
	}
	return this.repo.InsertFailedJob(item.ConvertToFailedJob(errorMessage))
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
