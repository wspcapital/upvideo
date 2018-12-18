package usr

import "time"

type User struct {
	Id                           int
	AccountId                    int
	Email                        string
	PasswordHash                 string
	APIKey                       string
	ForgotPasswordToken          string
	ForgotPasswordTokenExpiredAt *time.Time
}
