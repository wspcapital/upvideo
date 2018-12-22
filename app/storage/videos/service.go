package videos

import (
	"database/sql"
)

type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Video, error) {
	params.Limit = 1
	params.Offset = 0
	_video, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_video) == 0 {
		return nil, sql.ErrNoRows
	}
	return _video[0], nil
}

func (this *Service) FindAll(params Params) ([]*Video, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Video) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Video) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Video) error {
	return this.repo.Delete(item)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
