package aws

import (
	"bitbucket.org/marketingx/upvideo/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	"io"
	"net/url"
)

var (
	awsSession *session.Session
	awsParams  *config.AWSParams
)

type AWSEmail struct {
	From    string // From source email
	To      string // To destination email(s)
	Subject string // Subject text to send
	Text    string // Text is the text body representation
	HTML    string // HTMLBody is the HTML body representation
	ReplyTo string // Reply-To email(s)
}

func AWSInitSession(params *config.AWSParams) (err error) {
	awsParams = params
	awsConf := &aws.Config{
		Region:      aws.String(awsParams.Region),
		Credentials: credentials.NewStaticCredentials(awsParams.AccessKeyId, awsParams.SecretKey, ""),
	}

	awsSession, err = session.NewSession(awsConf)

	return
}

func UploadS3File(key string, reader io.Reader) (string, error) {
	uploader := s3manager.NewUploader(awsSession)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(awsParams.Bucket),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		fmt.Println("Failed to upload to AWS S3: " + err.Error())
		return "", err
	}

	fmt.Println("Successfully uploaded to", result.Location)

	location, err := url.QueryUnescape(result.Location)
	if err != nil {
		fmt.Println("Failed to decode location: " + err.Error())
		return "", err
	}

	return location, err
}

func SendEmail(emailData AWSEmail) *ses.SendEmailOutput {
	// start a new ses session
	svc := ses.New(awsSession)

	body := &ses.Body{}
	if emailData.Text != "" || emailData.HTML == "" {
		body.Text = &ses.Content{
			Data:    aws.String(emailData.Text), // Required
			Charset: aws.String("UTF-8"),
		}
	}
	if emailData.HTML != "" {
		body.Html = &ses.Content{
			Data:    aws.String(emailData.HTML), // Required
			Charset: aws.String("UTF-8"),
		}
	}

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			ToAddresses: []*string{
				aws.String(emailData.To), // Required
				// More values...
			},
		},
		Message: &ses.Message{ // Required
			Body: body, // Required
			Subject: &ses.Content{ // Required
				Data:    aws.String(emailData.Subject), // Required
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(emailData.From), // Required

		ReplyToAddresses: []*string{
			aws.String(emailData.ReplyTo), // Required
			// More values...
		},
	}

	// send email
	resp, err := svc.SendEmail(params)

	if err != nil {
		fmt.Println(err.Error())
	}
	return resp
}
