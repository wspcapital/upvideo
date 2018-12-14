package aws

import (
	"bitbucket.org/marketingx/upvideo/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

var (
	aws_session *session.Session
	params      *config.AWSParams
)

func AWSInitSession(conf config.Config) (err error) {
	params = &conf.AWS

	awsConf := &aws.Config{
		Region:      aws.String(params.Region),
		Credentials: credentials.NewStaticCredentials(params.AccessKeyId, params.SecretKey, ""),
	}

	aws_session, err = session.NewSession(awsConf)

	return
}

func UploadS3File(key string, reader io.Reader) error {
	uploader := s3manager.NewUploader(aws_session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("Failed to upload to AWS S3: " + err.Error())
		return err
	}

	fmt.Println("Successfully uploaded to", result.Location)

	return err
}
