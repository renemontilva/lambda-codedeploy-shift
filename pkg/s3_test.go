package handlers

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

func TestDownloadObject(t *testing.T) {
	t.Parallel()
	object, err := NewS3Object("lambda-codedeploy-20/",
		"config/dSr6pcN",
		S3ObjectWithAWSConfig(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://localhost:4566"),
			Credentials: credentials.NewStaticCredentials(
				"fakeAccessKey",
				"fakeSecretKey",
				"fakeToken"),
		}),
	)
	if err != nil {
		t.Errorf("Error returned from NewS3Object: %s", err)
	}

	err = object.download()
	if err != nil {
		t.Errorf("Error returned from DownloadConfigFile: %s", err)
	}
	assert.FileExists(t, object.LocalPath)
}

func TestExtractZipFile(t *testing.T) {
	object, err := NewS3Object("lambda-codedeploy-20/",
		"config/dSr6pcN",
		S3ObjectWithLocalPath("testdata/dSr6pcN"),
		S3ObjectWithUnZipPath("testdata/tmp"),
		S3ObjectWithAWSConfig(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://localhost:4566"),
			Credentials: credentials.NewStaticCredentials(
				"fakeAccessKey",
				"fakeSecretKey",
				"fakeToken"),
		}),
	)
	if err != nil {
		t.Errorf("Error returned from NewS3Object: %s", err)
	}
	err = object.extractZipFile()
	if err != nil {
		t.Errorf("Error returned from extractZipFile: %s", err)
	}
	assert.FileExists(t, "testdata/tmp/deployShift.yaml")
}
