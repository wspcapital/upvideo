package accounts

type Params struct {
	Id            int    `json:"id"`
	UserId        int    `json:"user_id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ChannelName   string `json:"channel_name"`
	ChannelUrl    string `json:"channel_url"`
	ClientId      string `json:"client_id"`
	ClientSecrets string `json:"client_secrets"`
	ClientSecretsRow string  `json:"client_secrets_row"`
	RequestToken  string `json:"request_token"`
	RequestTokenRow string   `json:"request_token_row"`
	AuthUrl       string `json:"auth_url"`
	OTPCode       string `json:"otpcode"`
	Note          string `json:"note"`
	OperationId   string `json:"operation_id"`
	Offset        uint64 `json:"offset"`
	Limit         uint64 `json:"limit"`
}
