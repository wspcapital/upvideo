package accounts

import (
	"database/sql"
)

type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Account, error) {
	params.Limit = 1
	params.Offset = 0
	_account, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(_account) == 0 {
		return nil, sql.ErrNoRows
	}
	return _account[0], nil
}

func (this *Service) FindAll(params Params) ([]*Account, error) {
	return this.repo.FindAll(params)
}

func (this *Service) FindByOperation(UserId int, OperationId string) (*Account, error) {
	return this.repo.FindByOperation(UserId, OperationId)
}

func (this *Service) Insert(item *Account) error {
	return this.repo.Insert(item)
}

func (this *Service) Update(item *Account) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Account) error {
	return this.repo.Delete(item)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
