package titles

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Title, err error) {
	query := sq.Select("Id", "UserId", "CampaignId", "Title", "Tags", "File", "TmpFile", "YoutubeId", "Posted", "Converted", "Pending", "FrameRate", "Resolution", "IpAddress").From("titles")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.CampaignId != 0 {
		query = query.Where("CampaignId = ?", params.CampaignId)
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
		title := &Title{}
		rows.Scan(
			&title.Id,
			&title.UserId,
			&title.CampaignId,
			&title.Title,
			&title.Tags,
			&title.File,
			&title.TmpFile,
			&title.YoutubeId,
			&title.Posted,
			&title.Converted,
			&title.Pending,
			&title.FrameRate,
			&title.Resolution,
			&title.IpAddress,
		)

		items = append(items, title)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Title) error {
	result, err := this.db.Exec("INSERT INTO titles(UserId, CampaignId, Title, Tags, File, TmpFile, YoutubeId, Posted, Converted, Pending, FrameRate, Resolution, IpAddress) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.UserId,
		item.CampaignId,
		item.Title,
		item.Tags,
		item.File,
		item.TmpFile,
		item.YoutubeId,
		item.Posted,
		item.Converted,
		item.Pending,
		item.FrameRate,
		item.Resolution,
		item.IpAddress,
	)

	if err != nil {
		fmt.Printf("SQL Insert err: \n%v\n", err)
		return err
	}

	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Title) error {
	_, err := this.db.Exec("UPDATE titles SET UserId=?, CampaignId=?, Title=?, Tags=?, File=?, TmpFile=?, YoutubeId=?, Posted=?, Converted=?, Pending=?, FrameRate=?, Resolution=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.CampaignId,
		item.Title,
		item.Tags,
		item.File,
		item.TmpFile,
		item.YoutubeId,
		item.Posted,
		item.Converted,
		item.Pending,
		item.FrameRate,
		item.Resolution,
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
