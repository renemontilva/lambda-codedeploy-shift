package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func initLogger() {

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// Handler is the lambda handler invoked by the `lambda.Start` function call
func Handler(cxt context.Context, event events.CodePipelineJobEvent) error {
	initLogger()
	environment := os.Getenv("ENVIRONMENT")
	// Get bucket and key from event
	logger.Info("Event received",
		slog.String("bucket", event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.BucketName),
		slog.String("key", event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.ObjectKey),
	)

	// Pipeline
	pipeline, err := NewPipeline(event.CodePipelineJob.ID)
	if err != nil {
		logger.Error("Error returned from NewPipeline",
			slog.String("Detail", err.Error()),
		)
		return err
	}

	// Create S3Object
	object, err := NewS3Object(
		fmt.Sprintf("%s/", event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.BucketName),
		event.CodePipelineJob.Data.InputArtifacts[0].Location.S3Location.ObjectKey,
		S3ObjectWithAWSConfig(&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				event.CodePipelineJob.Data.ArtifactCredentials.AccessKeyID,
				event.CodePipelineJob.Data.ArtifactCredentials.SecretAccessKey,
				event.CodePipelineJob.Data.ArtifactCredentials.SessionToken),
		}),
	)
	if err != nil {
		logger.Error("Error returned from NewS3Object",
			slog.String("Detail", err.Error()),
		)
		return err
	}
	// Get files from S3 and save to /tmp
	folderPath, err := object.Download()
	if err != nil {
		pipeline.JobFailed(err.Error())
		logger.Error("Error returned from Download",
			slog.String("Detail", err.Error()),
		)
		return err
	}

	// Create Deploy object
	deploy, err := NewDeploy(folderPath, DeployWithTaskDefFileName(fmt.Sprintf("taskdef_%s.json", environment)))
	if err != nil {
		pipeline.JobFailed(err.Error())
		logger.Error("error returned from NewDeploy",
			slog.String("Detail", err.Error()),
		)
		return fmt.Errorf("error returned from NewDeploy: %w", err)
	}
	// Preferred Strategy
	deploymentGroup, err := deploy.PreferredStrategy()
	if err != nil {
		pipeline.JobFailed(err.Error())
		logger.Error("Error returned from PreferredStrategy",
			slog.String("Detail", err.Error()),
		)
		return fmt.Errorf("error returned from PreferredStrategy: %w", err)
	}
	logger.Info("Preferred Strategy", slog.String("DeploymentGroup", deploymentGroup))
	// Deploy Strategy to CodeDeploy
	logger.Info("Deploying Preferred Strategy",
		slog.String("ApplicationName", deploy.DeployShift.ApplicationName),
		slog.String("DeploymentGroup", deploymentGroup),
		slog.String("deployShiftFileName", deploy.DeployShiftFileName),
		slog.String("taskDefFileName", deploy.TaskDefFileName),
		slog.String("appSpecFileName", deploy.AppSpecFileName),
	)
	deployID, err := deploy.DeployStrategy(event.CodePipelineJob.ID)
	if err != nil {
		pipeline.JobFailed(err.Error())
		logger.Error("Error returned from DeployStrategy",
			slog.String("Detail", err.Error()),
		)
		return err
	}
	logger.Info("Deployed Preferred Strategy",
		slog.String("DeploymentId", deployID),
	)
	// Send Success to CodePipeline
	pipeline.JobSucceeded()

	return nil
}
