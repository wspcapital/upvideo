package usr

type User struct {
	Id           int
	AccountId    int
	Email        string `validate:"required,email"`
	PasswordHash string
	APIKey       string
}
