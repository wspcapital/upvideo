package usr

import "time"

type UserSearchDto struct {
	Id                           string
	Email                        string
	PasswordHash                 string
	APIKey                       string
	ForgotPasswordToken          string
	ForgotPasswordTokenExpiredAt *time.Time
	Offset                       string
	Limit                        string
}
