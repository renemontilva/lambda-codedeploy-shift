package handlers

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Object struct {
	BucketName string
	ObjectKey  string
	LocalPath  string
	Config     *aws.Config
}

type option func(*S3Object) error

// NewS3Object returns a S3Object struct.
func NewS3Object(bucketName string, objectKey string, options ...option) (S3Object, error) {
	object := S3Object{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}
	for _, option := range options {
		err := option(&object)
		if err != nil {
			return S3Object{}, err
		}
	}
	return object, nil
}

// S3ObjectWithLocalPath returns a option function to set LocalPath field.
func S3ObjectWithLocalPath(localPath string) option {
	return func(object *S3Object) error {
		object.LocalPath = localPath
		return nil
	}
}

// S3ObjectWithAWSConfig returns a option function to set Config field.
func S3ObjectWithAWSConfig(config *aws.Config) option {
	return func(object *S3Object) error {
		object.Config = config
		return nil
	}

}

// Download downloads a S3 Object to a local file.
func (object *S3Object) Download() error {
	sess, err := session.NewSession(object.Config)
	downloader := s3manager.NewDownloader(sess)
	// Create a file to write the S3 Object contents to.
	configFile, err := os.Create(object.LocalPath)
	defer configFile.Close()

	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(configFile, &s3.GetObjectInput{
		Bucket: aws.String(object.BucketName),
		Key:    aws.String(object.ObjectKey),
	})
	if err != nil {
		return fmt.Errorf("Error downloading file: %w", err)
	}
	return nil
}
