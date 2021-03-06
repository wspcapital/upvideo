package accounts

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/kataras/iris/core/errors"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (this *Repository) FindByOperation(UserId int, OperationId string) (*Account, error) {
	query := sq.Select("Id", "UserId", "Username", "Password", "Channelname", "Channelurl", "ClientId", "Clientsecrets","Clientsecretsrow", "Requesttoken","Requesttokenrow", "AuthUrl", "OnetimeCode", "Note", "OperationId", "Created_at", "Updated_at").From("accounts")

	if UserId < 1 || OperationId == "" {
		return nil, errors.New("Wrong sql request parameters")
	}

	query = query.Where("UserId = ?", UserId).Where("OperationId = ?", OperationId)

	row := query.RunWith(this.db).QueryRow()

	_account := &Account{}
	err := row.Scan(
		&_account.Id,
		&_account.UserId,
		&_account.Username,
		&_account.Password,
		&_account.ChannelName,
		&_account.ChannelUrl,
		&_account.ClientId,
		&_account.ClientSecrets,
		&_account.ClientSecretsRow,
		&_account.RequestToken,
		&_account.RequestTokenRow,
		&_account.AuthUrl,
		&_account.OTPCode,
		&_account.Note,
		&_account.OperationId,
		&_account.Created_at,
		&_account.Updated_at,
	)
	log.Println(_account)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return _account, nil
}

func (this *Repository) FindAll(params Params) (items []*Account, err error) {
	query := sq.Select("Id", "UserId", "Username", "Password", "Channelname", "Channelurl", "ClientId", "Clientsecrets", "Clientsecretsrow", "Requesttoken", "Requesttokenrow", "AuthUrl", "OnetimeCode", "Note", "OperationId", "Created_at", "Updated_at").From("accounts")

	if params.UserId != 0 {
		query = query.Where("UserId = ?", params.UserId)
	}

	if params.ClientId != "" {
		query = query.Where("ClientId = ?", params.ClientId)
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
			&_account.Id,
			&_account.UserId,
			&_account.Username,
			&_account.Password,
			&_account.ChannelName,
			&_account.ChannelUrl,
			&_account.ClientId,
			&_account.ClientSecrets,
			&_account.ClientSecretsRow,
			&_account.RequestToken,
			&_account.RequestTokenRow,
			&_account.AuthUrl,
			&_account.OTPCode,
			&_account.Note,
			&_account.OperationId,
			&_account.Created_at,
			&_account.Updated_at,
		)
		items = append(items, _account)
	}
	rows.Close()
	return
}

func (this *Repository) Insert(item *Account) error {
	log.Println(item)
	result, err := this.db.Exec(`INSERT accounts SET 
                    UserId=?, 
                    Username=?, 
                    Password=?, 
                    Channelname=?, 
                    Channelurl=?, 
                    ClientId=?, 
                    Clientsecrets=?, 
                    Clientsecretsrow=?,
                    Requesttoken=?,
                    Requesttokenrow=?, 
                    AuthUrl=?, 
                    OnetimeCode=?, 
                    Note=?, 
                    OperationId=?`,
		item.UserId,
		item.Username,
		item.Password,
		item.ChannelName,
		item.ChannelUrl,
		item.ClientId,
		item.ClientSecrets,
		item.ClientSecretsRow,
		item.RequestToken,
		item.RequestTokenRow,
		item.AuthUrl,
		item.OTPCode,
		item.Note,
		item.OperationId,
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

func (this *Repository) Update(item *Account) error {
	_, err := this.db.Exec("UPDATE accounts SET UserId=?, Username=?, Password=?, Channelname=?, Channelurl=?, ClientId=?, Clientsecrets=?, Clientsecretsrow=?, Requesttoken=?, Requesttokenrow=?, AuthUrl=?, OnetimeCode=?, Note=?, OperationId=?, Updated_at=NOW() WHERE Id=?",
		item.UserId,
		item.Username,
		item.Password,
		item.ChannelName,
		item.ChannelUrl,
		item.ClientId,
		item.ClientSecrets,
		item.ClientSecretsRow,
		item.RequestToken,
		item.RequestTokenRow,
		item.AuthUrl,
		item.OTPCode,
		item.Note,
		item.OperationId,
		item.Id,
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
