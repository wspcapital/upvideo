package accounts

type Account struct {
	Id            int    `json:"id"`
	UserId        int    `json:"user_id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ChannelName   string `json:"channel_name"`
	ChannelUrl    string `json:"channel_url"`
	ClientSecrets string `json:"client_secrets"`
	RequestToken  string `json:"request_token"`
	AuthUrl       string `json:"auth_url"`
	OTPCode       string `json:"otpcode"`
	Note          string `json:"note"`
	OperationId   string `json:"operation_id"`
}
