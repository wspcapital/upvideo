package config

type EmailParams struct {
	From      string `json:"from"`
	ReplyTo   string `json:"reply"`
	Templates struct {
		RestorePasswordPath string `json:"restore_password"`
	} `json:"templates"`
}
