package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// File name used for update codedeploy
var deployShiftFile = "deployShift.yaml"

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(cxt context.Context, event events.CodePipelineJobEvent) error {
	// Get bucket and key from event
	object, err := NewS3Object(
		fmt.Sprintf("%s/", event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.BucketName),
		deployShiftFile,
		S3ObjectWithLocalPath(fmt.Sprintf("/tmp/%s", deployShiftFile)),
		S3ObjectWithAWSConfig(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://localstack:4566"),
			Credentials: credentials.NewStaticCredentials(
				event.CodePipelineJob.Data.ArtifactCredentials.AccessKeyID,
				event.CodePipelineJob.Data.ArtifactCredentials.SecretAccessKey,
				event.CodePipelineJob.Data.ArtifactCredentials.SessionToken),
		}),
	)
	if err != nil {
		return fmt.Errorf("Error returned from NewS3Object: %w", err)
	}
	// Get file from S3 and save to /tmp
	object.Download()
	// Read file from /tmp
	deployShift, err := NewDeployShift(object.LocalPath)
	if err != nil {
		return fmt.Errorf("Error returned from NewDeployShift: %w", err)
	}
	// Current Strategy
	deploymentGroup, err := deployShift.CurrentStrategy()
	if err != nil {
		return fmt.Errorf("Error returned from CurrentStrategy: %w", err)
	}
	// Update Strategy
	deployShift.DeployStrategy()
	log.Printf("ApplicationName: %s Current Strategy: %s", deployShift.ApplicationName, deploymentGroup)
	return nil
}
