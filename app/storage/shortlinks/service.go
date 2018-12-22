package shortlinks

import (
	"database/sql"
)

type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Shortlink, error) {
	params.Limit = 1
	params.Offset = 0
	_shortlink, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_shortlink) == 0 {
		return nil, sql.ErrNoRows
	}
	return _shortlink[0], nil
}

func (this *Service) FindAll(params Params) ([]*Shortlink, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Shortlink) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Shortlink) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Shortlink) error {
	return this.repo.Delete(item)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
