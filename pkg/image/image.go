package image

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	region       = "ap-northeast-2"
	s3BucketName = "erawgeragf"
)

var (
	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
)

func Upload(f *multipart.FileHeader, imageName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey,
			awsSecretKey,
			""),
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	file, _ := f.Open()
	if err != nil {
		fmt.Println("Error copying the file")
		return err
	}

	file.Seek(0, 0)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(imageName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func Download(imageName string) (*s3.GetObjectOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey,
			awsSecretKey,
			""),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
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
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey,
			awsSecretKey,
			""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(imageName),
	})
	if err != nil {
		return err
	}

	return nil
}
