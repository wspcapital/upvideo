package config

type AWSParams struct {
	AccessKeyId string `json:"ACCESS_KEY_ID"`
	SecretKey   string `json:"SECRET_KEY"`
	Region      string `json:"REGION"`
	Bucket      string `json:"bucket"`
}
