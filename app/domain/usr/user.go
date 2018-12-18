package usr

import "time"

type User struct {
	Id                           int        `json:"id"`
	AccountId                    int        `json:"account_id"`
	Email                        string     `json:"email"`
	PasswordHash                 string     `json:"password_hash"`
	APIKey                       string     `json:"api_key"`
	ForgotPasswordToken          string     `json:"forgot_password_token"`
	ForgotPasswordTokenExpiredAt *time.Time `json:"forgot_password_token_expired_at"`
}
