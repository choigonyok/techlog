package image

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	Region       = "ap-northeast-2"
	S3BucketName = "techlog-choigonyok"
)

func Upload(r io.Reader, id string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY,
			AWS_SECRET_KEY,
			""),
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(id),
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func Download(imageName string) (*s3.GetObjectOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY,
			AWS_SECRET_KEY,
			""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(imageName),
	})
	if err != nil {
		fmt.Println("ERROR:  ", err.Error())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded")
	return output, nil
}

func Remove(imageName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY,
			AWS_SECRET_KEY,
			""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(imageName),
	})
	if err != nil {
		return err
	}

	return nil
}
