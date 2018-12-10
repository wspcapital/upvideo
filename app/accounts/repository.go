package accounts

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindAll(params Params) (items []*Account, err error) {
	query := sq.Select("Id", "UserId", "Username", "Channelname").From("accounts")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.Username != "" {
		query = query.Where(sq.Eq{"Username": params.Username})
	}

	if params.ChannelName != "" {
		query = query.Where(sq.Eq{"Channelname": params.ChannelName})
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
		_account := &Account{}
		rows.Scan(	
					&_account.UserId,
					&_account.Username,
					&_account.Password,
					&_account.ChannelName,
					&_account.ChannelUrl,
					&_account.ClientSecrets,
					&_account.RequestToken,
					&_account.AuthUrl,
					&_account.OTPCode,
					&_account.Note,
				)
		log.Println(_account)
		items = append(items, _account)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Account) error {
	log.Println(item)
	result, err := this.db.Exec("INSERT INTO accounts(UserId, Username, Password, Channelname, Channelurl, Clientsecrets, Requesttoken, AuthUrl, OnetimeCode, Note, Created_at, Updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 
						item.UserId,
						item.Username,
						item.Password,
						item.ChannelName,
						item.ChannelUrl,
						item.ClientSecrets,
						item.RequestToken,
						item.AuthUrl,
						item.OTPCode,
						item.Note,
						int32(time.Now().Unix()),
						int32(time.Now().Unix()),
					)
	if err != nil {
		log.Println(err)
		return err
	}
	
	
	Id64, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	item.Id = int(Id64)
	return err
}

func (this *Repository) Update(item *Account) error {
	_, err := this.db.Exec("UPDATE accounts SET UserId=?, Username=?, Password=?, Channelname=?, Channelurl=?, Clientsecrets=?, Requesttoken=?, AuthUrl=?, OnetimeCode=?, Note=?, WHERE Id=?", 
						item.UserId,
						item.Username,
						item.Password,
						item.ChannelName,
						item.ChannelUrl,
						item.ClientSecrets,
						item.RequestToken,
						item.AuthUrl,
						item.OTPCode,
						item.Note,
					)
	return err
}

func (this *Repository) Delete(item *Account) error {
	_, err := this.db.Exec("DELETE FROM accounts WHERE Id=?", item.Id)
	return err
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
