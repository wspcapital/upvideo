package infrastructure

import (
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"database/sql"
	"encoding/json"
	"time"
)

type DbSession struct {
	db *sql.DB
}

func (this *DbSession) Cleanup(ttlMinutes int32) {
	this.db.Exec("delete from session where updated_at < ?", int32(time.Now().Unix())-ttlMinutes*60)
}

func (this *DbSession) FindById(id string) (item *session.Entity, err error) {
	var data string
	item = &session.Entity{}
	item.Data = make(map[string]string)
	err = this.db.QueryRow("select id, data from session where id = ?", id).Scan(&item.Id, &data)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(data), &item.Data)
	return
}

func (this *DbSession) Update(item *session.Entity) (err error) {
	var data []byte
	data, err = json.Marshal(item.Data)
	if err != nil {
		return
	}
	_, err = this.db.Exec("replace into session (id, data, updated_at) values (?, ?, ?)", item.Id, string(data), int32(time.Now().Unix()))
	return
}

func (this *DbSession) DeleteById(id string) error {
	_, err := this.db.Exec("DELETE FROM `session` WHERE id = ?", id)
	return err
}

func NewDbSession(db *sql.DB) *DbSession {
	return &DbSession{db: db}
}
