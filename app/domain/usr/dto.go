package usr

type UserSearchDto struct {
	Id           string
	Email        string
	PasswordHash string
	APIKey       string
	Offset       string
	Limit        string
}
