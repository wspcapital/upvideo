package shortlinks

import (
	"database/sql"
	"bitbucket.org/marketingx/upvideo/app/utils/checklinks"
	"fmt"
	_"strconv"
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

func (this *Service) CheckAllLinks() error {
	var _shortlink []*Shortlink

	fmt.Println("==> Start Check Short Links.")

	_shortlink, err := this.repo.FindAll(Params{
		Id: 0,
		Disabled: false,
		})

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, link := range _shortlink {
		if !link.Disabled && checklinks.CheckDisabledUrl(link.Url) {
			link.Disabled = true

			err = this.repo.Update(link)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// fmt.Println("This shortlink id : ", link.UniqId, "has beed disabled")
		}
	}

	return nil
}


func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
