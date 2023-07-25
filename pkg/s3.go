package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Object struct {
	BucketName string
	ObjectKey  string
	LocalPath  string
	UnZipPath  string
	Config     *aws.Config
}

type option func(*S3Object) error

// NewS3Object returns a S3Object struct.
func NewS3Object(bucketName string, objectKey string, options ...option) (S3Object, error) {
	_, file := path.Split(objectKey)

	object := S3Object{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		LocalPath:  fmt.Sprintf("/tmp/%s", file),
		UnZipPath:  "/tmp",
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

// S3ObjectWithUnZipPath returns a option function to set UnZipPath field.
func S3ObjectWithUnZipPath(unZipPath string) option {
	return func(object *S3Object) error {
		object.UnZipPath = unZipPath
		return nil
	}
}

// Download downloads a S3 Zip Object to a local file and extracts it.
func (object *S3Object) Download() (string, error) {
	err := object.download()
	if err != nil {
		return "", err
	}
	err = object.extractZipFile()
	if err != nil {
		return "", err
	}
	return object.UnZipPath, nil
}

// Download downloads a S3 Object to a local file.
func (object *S3Object) download() error {
	sess, err := session.NewSession(object.Config)
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}
	downloader := s3manager.NewDownloader(sess)
	// Create a file to write the S3 Object contents to.
	configFile, err := os.Create(object.LocalPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer configFile.Close()

	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(configFile, &s3.GetObjectInput{
		Bucket: aws.String(object.BucketName),
		Key:    aws.String(object.ObjectKey),
	})
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}
	return nil
}

// extractZipFile extracts a zip file.
func (object *S3Object) extractZipFile() error {
	openZipFile, err := zip.OpenReader(object.LocalPath)
	if err != nil {
		return err
	}
	defer openZipFile.Close()
	for _, file := range openZipFile.File {
		filePath := path.Join(object.UnZipPath, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}
	return nil
}
