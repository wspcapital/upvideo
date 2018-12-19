package campaigns

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) fetch(query sq.SelectBuilder) ([]*Campaign, error) {
	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*Campaign, 0)

	for rows.Next() {
		item := new(Campaign)
		err = rows.Scan(
			&item.Id,
			&item.UserId,
			&item.AccountId,
			&item.VideoId,
			&item.Title,
			&item.TotalTitles,
			&item.CompleteTitles,
			&item.TitlesGenerated,
			&item.IpAddress,
			&item.DateStart_at,
			&item.DateComplete_at,
			&item.Created_at,
			&item.Updated_at,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (this *Repository) FindAll(params Params) (result []*Campaign, err error) {
	query := sq.Select(
		"Id",
		"UserId",
		"AccountId",
		"VideoId",
		"Title",
		"TotalTitles",
		"CompleteTitles",
		"TitlesGenerated",
		"IpAddress",
		"DateStart_at",
		"DateComplete_at",
		"Created_at",
		"Updated_at").From("campaigns")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.AccountId != 0 {
		query = query.Where("AccountId = ?", params.AccountId)
	}

	if params.VideoId != 0 {
		query = query.Where("VideoId = ?", params.VideoId)
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

	query.OrderBy("`sort` asc")

	result, err = this.fetch(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *Repository) Insert(item *Campaign) error {
	log.Println(item)
	result, err := this.db.Exec(
		"INSERT campaigns SET UserId=?, AccountId=?, VideoId=?, Title=?, IpAddress=?",
		item.UserId,
		item.AccountId,
		item.VideoId,
		item.Title,
		item.IpAddress,
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

func (this *Repository) Update(item *Campaign) error {
	_, err := this.db.Exec("UPDATE campaigns SET UserId=?, AccountId=?, VideoId=?, Title=?, TotalTitles=?, CompleteTitles=?, TitlesGenerated=?, DateStart_at=?, DateComplete_at=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.AccountId,
		item.VideoId,
		item.Title,
		item.TotalTitles,
		item.CompleteTitles,
		item.TitlesGenerated,
		item.DateStart_at,
		item.DateComplete_at,
		item.Id,
	)
	return err
}

func (this *Repository) Delete(item *Campaign) error {
	_, err := this.db.Exec("DELETE FROM campaigns WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
