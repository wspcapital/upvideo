package email

import (
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/aws"
	"bitbucket.org/marketingx/upvideo/config"
	"bytes"
	"fmt"
	"html/template"
)

type Service interface {
	SendRestorePasswordEmail(*usr.User) error
}

type service struct {
	emailParams *config.EmailParams
}

func (this *service) SendRestorePasswordEmail(user *usr.User) error {
	htmlBody, err := this.parseTemplate(this.emailParams.Templates.RestorePasswordPath, struct {
		Email string
		Token string
	}{
		Email: user.Email,
		Token: user.ForgotPasswordToken,
	})

	if err != nil {
		fmt.Println("\nSendRestorePasswordEmail -> template error:", err)
		return err
	}

	aws.SendEmail(aws.AWSEmail{
		From:    this.emailParams.From,
		ReplyTo: this.emailParams.ReplyTo,
		To:      user.Email,
		Subject: "Restore password request",
		Text:    "",
		HTML:    htmlBody,
	})

	return err
}

func (this *service) parseTemplate(tplFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(tplFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func NewService(emailParams *config.EmailParams) Service {
	return &service{
		emailParams: emailParams,
	}
}
