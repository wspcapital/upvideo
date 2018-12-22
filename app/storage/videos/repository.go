package videos

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Video, err error) {
	query := sq.Select("Id", "UserId", "AccountId", "Title", "Description", "Tags", "Category", "Language", "File", "TmpFile", "Playlist", "Title_gen", "IpAddress").From("videos")

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

	query.OrderBy("`sort` asc")

	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		return
	}
	for rows.Next() {
		video := &Video{}
		rows.Scan(
			&video.Id,
			&video.UserId,
			&video.AccountId,
			&video.Title,
			&video.Description,
			&video.Tags,
			&video.Category,
			&video.Language,
			&video.File,
			&video.TmpFile,
			&video.Playlist,
			&video.TitleGen,
			&video.IpAddress,
		)
		items = append(items, video)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Video) error {
	result, err := this.db.Exec("INSERT INTO videos(UserId, AccountId, Title, Description, Tags, Category, `Language`, File, TmpFile, Playlist, IpAddress, Created_at, Updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.UserId,
		item.AccountId,
		item.Title,
		item.Description,
		item.Tags,
		item.Category,
		item.Language,
		item.File,
		item.TmpFile,
		item.Playlist,
		item.IpAddress,
		int32(time.Now().Unix()),
		int32(time.Now().Unix()),
	)

	if err != nil {
		fmt.Printf("SQL Insert err: \n%v\n", err)
		return err
	}

	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Video) error {
	_, err := this.db.Exec("UPDATE videos SET UserId=?, AccountId=?, Title=?, Description=?, Tags=?, Category=?, `Language`=?, File=?, TmpFile=?, Playlist=?, Title_gen=?, IpAddress=? WHERE Id=?", item.UserId, item.AccountId, item.Title, item.Description, item.Tags, item.Category, item.Language, item.File, item.TmpFile, item.Playlist, item.TitleGen, item.IpAddress, item.Id)
	return err
}

func (this *Repository) Delete(item *Video) error {
	_, err := this.db.Exec("DELETE FROM videos WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
