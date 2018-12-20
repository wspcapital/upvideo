package titles

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/guregu/null"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Title, err error) {
	query := sq.Select("Id", "UserId", "CampaignId", "Title", "Tags", "File", "TmpFile", "YoutubeId", "YoutubeUrl", "Posted", "Converted", "Pending", "FrameRate", "Resolution", "IpAddress", "Created_at", "Updated_at").From("titles")

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
			&title.YoutubeUrl,
			&title.Posted,
			&title.Converted,
			&title.Pending,
			&title.FrameRate,
			&title.Resolution,
			&title.IpAddress,
			&title.CreatedAt,
			&title.UpdatedAt,
		)

		items = append(items, title)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Title) error {
	type maxFrameRate struct {
		FrameRate  null.Int
		Resolution null.Int
	}

	mfr := &maxFrameRate{}

	err := this.db.QueryRow("SELECT MAX(FrameRate) AS maxFrameRate, Resolution AS maxResolution FROM titles WHERE Resolution in (SELECT MAX(Resolution) FROM titles WHERE CampaignId=?) GROUP BY Resolution LIMIT 1", item.CampaignId).
		Scan(&mfr.FrameRate, &mfr.Resolution)
	if err != nil {
		return err
	}

	if mfr.Resolution.Valid && mfr.FrameRate.Valid {
		item.Resolution = int(mfr.Resolution.Int64)
		item.FrameRate = int(mfr.FrameRate.Int64)
	}

	var result sql.Result
	insert := false

	for !insert {
		item.SetNextFrameRate()

		result, err = this.db.Exec(`INSERT titles SET 
                  UserId=?, 
                  CampaignId=?, 
                  Title=?, 
                  Tags=?, 
                  File=?, 
                  TmpFile=?, 
                  FrameRate=?, 
                  Resolution=?, 
                  IpAddress=?`,
			item.UserId,
			item.CampaignId,
			item.Title,
			item.Tags,
			item.File,
			item.TmpFile,
			item.FrameRate,
			item.Resolution,
			item.IpAddress,
		)

		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number != 1062 { // Duplicate error
				return err
			} else { // title can be duplicated too
				has, err2 := this.Has(item)
				if err2 != nil {
					return err2
				}
				if has {
					return err
				}
			}
		} else {
			insert = true
		}
	}

	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *Repository) Has(item *Title) (bool, error) {
	var count int
	err := this.db.QueryRow("SELECT COUNT(id) FROM titles WHERE CampaignId=? AND Title=? GROUP BY Resolution", item.CampaignId, item.Title).
		Scan(&count)
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (this *Repository) Update(item *Title) error {
	_, err := this.db.Exec("UPDATE titles SET UserId=?, CampaignId=?, Title=?, Tags=?, File=?, TmpFile=?, YoutubeId=?, YoutubeUrl=?, Posted=?, Converted=?, Pending=?, FrameRate=?, Resolution=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.CampaignId,
		item.Title,
		item.Tags,
		item.File,
		item.TmpFile,
		item.YoutubeId,
		item.YoutubeUrl,
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
