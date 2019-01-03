package topusers

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*TopUser, err error) {

	rows, err := this.db.Query("SELECT `user`.Id as id, " +
	"`user`.Email as author_name, " +
		"COUNT(DISTINCT accounts.id) AS new_channels, " +
		"MAX(accounts.Updated_at) AS last_activity, " +
		"COUNT(DISTINCT campaigns.id) AS new_videos " +
		"FROM accounts " +
	"LEFT OUTER JOIN `user` ON (accounts.UserId = `user`.id) " +
	"LEFT OUTER JOIN videos ON (accounts.id = videos.AccountId) " +
	"LEFT OUTER JOIN campaigns ON (videos.id = campaigns.VideoId) " +
	"WHERE (accounts.Deleted = FALSE " +
	"AND accounts.Created_at BETWEEN ? AND ?) " +
	"GROUP BY user.Email, " +
		"user.Id " +
		"ORDER BY user.Email ASC ",
		params.GreaterThan,
		params.LessThan,
	)
	if err != nil {
		log.Println(`Error FindAll`, err)
		return
	}



	for rows.Next() {
		topuser := &TopUser{}
		rows.Scan(
			&topuser.Id,
			&topuser.Email,
			&topuser.NewChannelsCount,
			&topuser.LastActivity,
			&topuser.NewVideosCount,
		)

		items = append(items, topuser)
	}
	rows.Close()
	return
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
