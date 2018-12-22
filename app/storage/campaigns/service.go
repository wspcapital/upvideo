package campaigns

import (
	"database/sql"
)

type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Campaign, error) {
	params.Limit = 1
	params.Offset = 0
	result, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}
	return result[0], nil
}

func (this *Service) FindAll(params Params) ([]*Campaign, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Campaign) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Campaign) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Campaign) error {
	return this.repo.Delete(item)
}

func (this *Service) CountTotalTitles(item *Campaign) error {
	return this.repo.CountTotalTitles(item)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
