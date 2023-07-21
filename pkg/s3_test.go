package handlers

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func TestDownloadObject(t *testing.T) {
	t.Parallel()
	object, err := NewS3Object("codepipeline-us-east-1/",
		"deployShift.yaml",
		S3ObjectWithLocalPath("/tmp/deployShift.yaml"),
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

	err = object.Download()
	if err != nil {
		t.Errorf("Error returned from DownloadConfigFile: %s", err)
	}
	_, err = os.Stat(object.LocalPath)
	if err != nil {
		t.Errorf("Error checking file: %s", err)
	}
}
