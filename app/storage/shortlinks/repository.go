package shortlinks

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) fetch(query sq.SelectBuilder) ([]*Shortlink, error) {
	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*Shortlink, 0)

	for rows.Next() {
		item := new(Shortlink)
		err = rows.Scan(
			&item.Id,
			&item.UserId,
			&item.UniqId,
			&item.Url,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (this *Repository) FindAll(params Params) (result []*Shortlink, err error) {
	query := sq.Select(
		"Id",
		"UserId",
		"UniqId",
		"Url",
		"Created_at",
		"Updated_at").From("shortlinks")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.UniqId != "" {
		query = query.Where("UniqId = ?", params.UniqId)
	}

	if params.Url != "" {
		query = query.Where(sq.Eq{"Url": params.Url})
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

	query.OrderBy("`sort` asc")

	result, err = this.fetch(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *Repository) Insert(item *Shortlink) error {
	log.Println(item)
	result, err := this.db.Exec(
		"INSERT shortlinks SET UserId=?, UniqId=?, Url=?",
		item.UserId,
		item.UniqId,
		item.Url,
	)
	if err != nil {
		fmt.Printf("SQL Insert err: \n%v\n", err)
		return err
	}

	Id64, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Shortlink) error {
	_, err := this.db.Exec("UPDATE shortlinks SET UserId=?, UniqId=?, Url=?, Counter=?, Created_at=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.UniqId,
		item.Url,
		item.Counter,
		item.CreatedAt,
		item.Id,
	)
	return err
}

func (this *Repository) Delete(item *Shortlink) error {
	_, err := this.db.Exec("DELETE FROM shortlinks WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
