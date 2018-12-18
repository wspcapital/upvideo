package session

import "bitbucket.org/marketingx/upvideo/app/cookie"

type Service interface {
	FindById(id string) (*Entity, error)
	Update(item *Entity) error
	Create() *Entity
	DeleteById(id string) error
}

type service struct {
	repo Repository
}

func (this *service) FindById(id string) (*Entity, error) {
	return this.repo.FindById(id)
}

func (this *service) Update(item *Entity) error {
	return this.repo.Update(item)
}

func (this *service) Create() *Entity {
	key, err := cookie.Generate(128)
	if err != nil {
		return nil
	}
	item := &Entity{Id: key, Data: make(map[string]string)}
	err = this.repo.Update(item)
	if err != nil {
		return nil
	}
	return item
}

func (this *service) DeleteById(id string) error {
	return this.repo.DeleteById(id)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
