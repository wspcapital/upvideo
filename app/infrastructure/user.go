package infrastructure

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"strconv"
)

type UserRepository struct {
	db *sql.DB
}

func (this *UserRepository) FindAll(dto usr.UserSearchDto) (items []*usr.User, err error) {
	var offset, limit uint64
	query := sq.Select("Id", "Email", "PasswordHash", "APIKey").From("user")
	if dto.Id != "" {
		query = query.Where(sq.Eq{"Id": dto.Id})
		query = query.Limit(1)
		query = query.Offset(0)
	}

	if dto.Email != "" {
		query = query.Where(sq.Eq{"Email": dto.Email})
		query = query.Limit(1)
		query = query.Offset(0)
	}

	if dto.PasswordHash != "" {
		query = query.Where(sq.Eq{"PasswordHash": dto.PasswordHash})
		query = query.Limit(1)
		query = query.Offset(0)
	}

	if dto.APIKey != "" {
		query = query.Where(sq.Eq{"APIKey": dto.APIKey})
		query = query.Limit(1)
		query = query.Offset(0)
	}

	if dto.Limit != "" {
		limit, err = strconv.ParseUint(dto.Limit, 10, 64)
		if err != nil {
			return
		}
		query = query.Limit(limit)
	}

	if dto.Offset != "" {
		offset, err = strconv.ParseUint(dto.Offset, 10, 64)
		if err != nil {
			return
		}
		query = query.Offset(offset)
	}

	rows, err := query.RunWith(this.db).Query()
	if err != nil {
		return
	}
	for rows.Next() {
		user := &usr.User{}
		rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.APIKey)
		items = append(items, user)
	}
	rows.Close()
	return
}

func (this *UserRepository) Insert(item *usr.User) error {
	result, err := this.db.Exec("INSERT INTO user(Email, PasswordHash, APIKey) VALUES(?, ?, ?)", item.Email, item.PasswordHash, item.APIKey)
	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *UserRepository) Update(item *usr.User) error {
	_, err := this.db.Exec("update user set Email=?, PasswordHash=?, APIKey=? where Id=?", item.Email, item.PasswordHash, item.Id, item.APIKey)
	return err
}

func (this *UserRepository) Delete(item *usr.User) error {
	_, err := this.db.Exec("delete from user where Id=?", item.Id)
	return err
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
