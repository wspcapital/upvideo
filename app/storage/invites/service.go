package invites

import (
	"bitbucket.org/marketingx/upvideo/app/cookie"
	"errors"
)


type Service struct {
	repo *Repository
}

func (this *Service) FindOne(params Params) (*Invite, error) {
	params.Limit = 1
	params.Offset = 0
	invites, err := this.repo.FindAll(params)
	if err != nil {
		return nil, err
	}
	if len(invites) == 0 {
		return nil, errors.New("found no matching invites")
	}
	return invites[0], nil
}

func (this *Service) FindAll(params Params) ([]*Invite, error) {
	return this.repo.FindAll(params)
}

func (this *Service) Insert(item *Invite) error {
	code, err := cookie.Generate(96)
	if err != nil {
		return err
	}
	item.Code = code
	return this.repo.Insert(item)
}


func (this *Service) Update(item *Invite) error {
	return this.repo.Update(item)
}

func (this *Service) Delete(item *Invite) error {
	return this.repo.Delete(item)
}

func (this *Service) CheckInvite(code string) error {
	if code == "" {
		return errors.New("no code")
	}
	// FIXME should we check if a user is allowed to invite?
	_, err := this.FindOne(Params{
		Code: code,
	})
	if err != nil {
		return err
	}
	return nil
}

func (this *Service) ClearInvite(code string) error {
	if code == "" {
		return errors.New("no code")
	}
	invite, err := this.FindOne(Params{
		Code: code,
	})
	if err != nil {
		return err
	}
	return this.repo.Delete(invite)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
