package invites

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Invite, err error) {
	query := sq.Select("Id", "UserId", "Title", "Code").From("invite")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.Title != "" {
		query = query.Where(sq.Eq{"Title": params.Title})
	}

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
		query = query.Offset(params.Offset)
	}

	if params.Id != 0 {
		query = query.Where("Id = ?", params.Id)
		query = query.Limit(1)
		query = query.Offset(0)
	}

	if params.Code != "" {
		query = query.Where("Code = ?", params.Code)
		query = query.Limit(1)
		query = query.Offset(0)
	}

	query.OrderBy("`sort` asc")

	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		return
	}
	for rows.Next() {
		invite := &Invite{}
		rows.Scan(&invite.Id, &invite.UserId, &invite.Title, &invite.Code)
		items = append(items, invite)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Invite) error {
	result, err := this.db.Exec("INSERT INTO invite(UserId, Title, Code) VALUES(?, ?, ?)", item.UserId, item.Title, item.Code)
	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Invite) error {
	_, err := this.db.Exec("UPDATE invite SET UserId=?, Title=?, Code=? WHERE Id=?", item.UserId, item.Title, item.Code, item.Id)
	return err
}

func (this *Repository) Delete(item *Invite) error {
	_, err := this.db.Exec("DELETE FROM invite WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
