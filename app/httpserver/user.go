package httpserver

import (
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"strconv"
)

type UserRepository struct {
	db *sql.DB
}

func (this *UserRepository) FindAll(dto usr.UserSearchDto) (items []*usr.User, err error) {
	var offset, limit uint64
	query := sq.Select("Id", "Email", "PasswordHash", "APIKey", "AccountId").From("user")
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
		rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.APIKey, &user.AccountId)
		items = append(items, user)
	}
	rows.Close()
	return
}

func (this *UserRepository) Insert(item *usr.User) error {
	result, err := this.db.Exec("INSERT INTO user(Email, PasswordHash, APIKey, AccountId) VALUES(?, ?, ?, ?)", item.Email, item.PasswordHash, item.APIKey, item.AccountId)
	Id64, err := result.LastInsertId()
	item.Id = int(Id64)
	return err
}

func (this *UserRepository) Update(item *usr.User) error {
	_, err := this.db.Exec("update user set Email=?, PasswordHash=?, APIKey=?, AccountId=? where Id=?", item.Email, item.PasswordHash, item.APIKey, item.AccountId, item.Id)
	return err
}

func (this *UserRepository) Delete(item *usr.User) error {
	_, err := this.db.Exec("delete from user where Id=?", item.Id)
	return err
}

func (this *UserRepository) FindByForgotPasswordToken(dto *usr.UserSearchDto) (*usr.User, error) {
	query := sq.Select("Id", "Email", "PasswordHash", "APIKey", "AccountId").From("user")
	query = query.
		Where(sq.Eq{"Email": dto.Email}).
		Where(sq.Eq{"ForgotPasswordToken": dto.ForgotPasswordToken}).
		Where("ForgotPasswordToken IS NOT NULL").
		Where("ForgotPasswordTokenExpiredAt IS NOT NULL AND ForgotPasswordTokenExpiredAt > NOW()").
		Limit(1).Offset(0)

	user := &usr.User{}
	err := query.RunWith(this.db).QueryRow().Scan(&user.Id, &user.Email, &user.PasswordHash, &user.APIKey, &user.AccountId)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (this *UserRepository) SetForgotPasswordToken(item *usr.User) error {
	_, err := this.db.Exec("update user set ForgotPasswordToken=?, ForgotPasswordTokenExpiredAt=DATE_ADD(NOW(), INTERVAL 1 DAY) where Id=?", item.ForgotPasswordToken, item.Id)
	return err
}

func (this *UserRepository) RemoveForgotPasswordToken(item *usr.User) error {
	_, err := this.db.Exec("update user set ForgotPasswordToken=NULL , ForgotPasswordTokenExpiredAt=NULL where Id=?", item.Id)
	return err
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
