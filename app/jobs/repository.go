package jobs

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Job, err error) {
	query := sq.Select("Id", "UserId", "RelatedId", "Type", "ProcessId", "Progress").From("jobs")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.RelatedId != 0 {
		query = query.Where("RelatedId = ?", params.RelatedId)
	}

	if params.Type != "" {
		query = query.Where(sq.Eq{"Type": params.Type})
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
		_job := &Job{}
		rows.Scan(
			&_job.Id,
			&_job.UserId,
			&_job.RelatedId,
			&_job.Type,
			&_job.ProcessId,
			&_job.Progress,
		)
		items = append(items, _job)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Job) error {
	result, err := this.db.Exec("INSERT INTO jobs(UserId, RelatedId, `Type`, ProcessId, Progress) VALUES(?, ?, ?, ?, ?)",
		item.UserId,
		item.RelatedId,
		item.Type,
		item.ProcessId,
		item.Progress,
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

func (this *Repository) Update(item *Job) error {
	_, err := this.db.Exec("UPDATE jobs SET UserId=?, RelatedId=?, `Type`=?, ProcessId=?, Progress=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.RelatedId,
		item.Type,
		item.ProcessId,
		item.Progress,
		item.Id,
	)
	return err
}

func (this *Repository) Delete(item *Job) error {
	_, err := this.db.Exec("DELETE FROM jobs WHERE Id=?", item.Id)
	return err
}

func (this *Repository) FindAllFailedJobs(params Params) (items []*FailedJob, err error) {
	query := sq.Select("Id", "UserId", "RelatedId", "Type", "ProcessId", "Error").From("failed_jobs")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.RelatedId != 0 {
		query = query.Where("RelatedId = ?", params.RelatedId)
	}

	if params.Type != "" {
		query = query.Where(sq.Eq{"Type": params.Type})
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
		_job := &FailedJob{}
		rows.Scan(
			&_job.Id,
			&_job.UserId,
			&_job.RelatedId,
			&_job.Type,
			&_job.ProcessId,
			&_job.Error,
		)
		items = append(items, _job)
	}
	rows.Close()
	return
}

func (this *Repository) InsertFailedJob(item *FailedJob) error {
	result, err := this.db.Exec("INSERT INTO failed_jobs(UserId, RelatedId, `Type`, ProcessId, `Error`) VALUES(?, ?, ?, ?, ?)",
		item.UserId,
		item.RelatedId,
		item.Type,
		item.ProcessId,
		item.Error,
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

func (this *Repository) DeleteFailedJob(item *FailedJob) error {
	_, err := this.db.Exec("DELETE FROM failed_jobs WHERE Id=?", item.Id)
	return err
}

func (this *Repository) DeleteFailedJobByType(item *FailedJob) error {
	_, err := this.db.Exec("DELETE FROM failed_jobs WHERE RelatedId=? AND `Type`=?", item.RelatedId, item.Type)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
