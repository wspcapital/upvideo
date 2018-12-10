package titles

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Title, err error) {
	query := sq.Select("Id", "UserId", "VideoId", "Title").From("titles")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
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

	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		log.Println(`Error FindAll`, err)
		return
	}
	for rows.Next() {
		titlee := &Title{}
		rows.Scan(
			&titlee.UserId,
			&titlee.VideoId,
			&titlee.Title,
			&titlee.Tags,
			&titlee.File,
			&titlee.Posted,
			&titlee.Converted,
			&titlee.Pending,
			&titlee.IpAddress,
		)

		log.Println(titlee)

		items = append(items, titlee)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Title) error {
	result, err := this.db.Exec("INSERT INTO titles(UserId, VideoId, Title, Tags, File, Posted, Converted, Pending, IpAddress) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)",
					item.UserId,
					item.VideoId,
					item.Title,
					item.Tags,
					item.File,
					item.Posted,
					item.Converted,
					item.Pending,
					item.IpAddress,
				)
	log.Println(err)
	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Title) error {
	_, err := this.db.Exec("UPDATE titles SET UserId=?, VideoId=?, Title=?, Tags=?, File=?, Posted=?, Converted=?, Pending=?, IpAddress=? WHERE Id=?",
					item.UserId,
					item.VideoId,
					item.Title,
					item.Tags,
					item.File,
					item.Posted,
					item.Converted,
					item.Pending,
					item.IpAddress,
					item.Id,
				)
	return err
}

func (this *Repository) Delete(item *Title) error {
	_, err := this.db.Exec("DELETE FROM titles WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
