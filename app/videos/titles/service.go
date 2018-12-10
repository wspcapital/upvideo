package titles

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Title, error) {
	params.Limit = 1
	params.Offset = 0
	_title, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_title) == 0 {
		return nil, errors.New("found no matching title")
	}
	return _title[0], nil
}

func (this *Service) FindAll(params Params) ([]*Title, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Title) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Title) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Title) error {
	return this.repo.Delete(item)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
