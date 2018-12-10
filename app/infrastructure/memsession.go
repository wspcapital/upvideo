package infrastructure

import (
	"errors"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"sync"
)

type MemSession struct {
	items map[string]*session.Entity
}

func (this *MemSession) Cleanup(ttlMinutes int32) {
}

func (this *MemSession) FindById(id string) (item *session.Entity, err error) {
	var found bool
	item, found = this.items[id]
	if !found {
		err = errors.New("Session not found")
	}
	return
}

func (this *MemSession) Update(item *session.Entity) error {
	this.items[item.Id] = item
	return nil
}

var sessionContainer *MemSession
var sessionOnce sync.Once

func GetMemSession() *MemSession {
	sessionOnce.Do(func() {
		sessionContainer = &MemSession{}
		sessionContainer.items = make(map[string]*session.Entity)
	})
	return sessionContainer
}
