package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
)

type Deploy struct {
	AppSpecFileName     string
	TaskDefFileName     string
	DeployShiftFileName string
	ImageDetailFileName string
	FolderPath          string
	DeployShift         DeployShift
	Config              *aws.Config
	Client              *codedeploy.CodeDeploy
}
type DeployShift struct {
	ApplicationName   string           `yaml:"applicationName"`
	DeployStrategies  []DeployStrategy `yaml:"deployStrategies"`
	PreferredStrategy string           `yaml:"preferredStrategy"`
}

type DeployStrategy struct {
	Name            string `yaml:"name"`
	DeploymentGroup string `yaml:"deploymentGroup"`
}

type optionDeploy func(*Deploy) error

// NewDeploy returns a Deploy struct.
func NewDeploy(folderPath string, options ...optionDeploy) (Deploy, error) {
	deploy := Deploy{
		AppSpecFileName:     "appspec.yaml",
		TaskDefFileName:     "taskdef.json",
		DeployShiftFileName: "deployShift.yaml",
		ImageDetailFileName: "imageDetail.json",
		FolderPath:          folderPath,
		Config:              &aws.Config{},
	}

	for _, option := range options {
		err := option(&deploy)
		if err != nil {
			return Deploy{}, err
		}
	}

	err := convertYAMLToType(filepath.Join(folderPath, deploy.DeployShiftFileName), &deploy.DeployShift)
	if err != nil {
		return Deploy{}, fmt.Errorf("error returned from NewDeploy: %w", err)
	}

	deploy.Client, err = deploy.connectCodeDeploy()
	if err != nil {
		return Deploy{}, fmt.Errorf("error connecting to CodeDeploy in NewDeploy function: %w", err)
	}

	return deploy, nil
}

// DeployWithBuildSpecFileName returns a option function to set AppSpecFileName field.
func DeployWithAppSpecFileName(appSpecFileName string) optionDeploy {
	return func(deploy *Deploy) error {
		deploy.AppSpecFileName = appSpecFileName
		return nil
	}
}

// DeployWithTaskDefFileName returns a option function to set TaskDefFileName field.
func DeployWithTaskDefFileName(taskDefFileName string) optionDeploy {
	return func(deploy *Deploy) error {
		deploy.TaskDefFileName = taskDefFileName
		return nil
	}
}

// DeployWithDeployShiftFileName returns a option function to set DeployShiftFileName field.
func DeployWithDeployShiftFileName(deployShiftFileName string) optionDeploy {
	return func(deploy *Deploy) error {
		deploy.DeployShiftFileName = deployShiftFileName
		return nil
	}
}

// DeployWithAWSConfig returns a option function to set Config field.
func DeployWithAWSConfig(config *aws.Config) optionDeploy {
	return func(deploy *Deploy) error {
		deploy.Config = config
		return nil
	}
}

// DeployStrategy triggers a deployment to CodeDeploy.
func (deploy *Deploy) DeployStrategy(jobID string) (string, error) {

	deploymentGroup, err := deploy.PreferredStrategy()
	if err != nil {
		return "", fmt.Errorf("error getting deployment group in DeployStrategy function: %w", err)
	}

	// ImageDetail
	imageDetail, err := NewImageDetail(deploy.FolderPath, deploy.ImageDetailFileName)
	if err != nil {
		return "", fmt.Errorf("error getting image detail in DeployStrategy function: %w", err)
	}

	// Taskdef
	taskDef, err := NewTaskDef(deploy.FolderPath, deploy.TaskDefFileName, WithImageURI(imageDetail.ImageURI))
	if err != nil {
		return "", err
	}
	err = taskDef.NewRevision(deploy.Config)
	if err != nil {
		return "", fmt.Errorf("error creating TaskDef new revision: %w", err)
	}

	// AppSpec
	appSpec, err := NewAppSpec(deploy.FolderPath, deploy.AppSpecFileName, WithTaskDefinition(*taskDef.TaskDefinitionArn))
	if err != nil {
		return "", fmt.Errorf("error getting appSpec in DeployStrategy function: %w", err)
	}
	appSpecFile, err := appSpec.String()
	if err != nil {
		return "", fmt.Errorf("error getting appSpec.String in DeployStrategy function: %w", err)
	}

	// Deploy to CodeDeploy with AppSpec
	deployResult, err := deploy.createDeployment(deploymentGroup, appSpecFile)
	if err != nil {
		return "", fmt.Errorf("error creating deployment in DeployStrategy function: %w", err)
	}
	return *deployResult.DeploymentId, nil
}

// findDeploymentGroup returns a deploymentGroup name defined in PreferredStrategy field of deployShift.yaml.
func (deploy *Deploy) findDeploymentGroup() (string, error) {
	for _, deployStrategy := range deploy.DeployShift.DeployStrategies {
		if deployStrategy.Name == deploy.DeployShift.PreferredStrategy {
			return deployStrategy.DeploymentGroup, nil
		}
	}
	return "", fmt.Errorf("deployment group not found in deployShift.yaml")
}

// PreferredStrategy returns a deploymentGroup name defined in ActiveStrategy field of deployShift.yaml.
func (deploy *Deploy) PreferredStrategy() (string, error) {
	deploymentGroup, err := deploy.findDeploymentGroup()
	if err != nil {
		return "", fmt.Errorf("error getting deployment group in PreferredStrategy function: %w", err)
	}
	return deploymentGroup, nil
}

func (deploy *Deploy) createDeployment(deploymentGroup string, appSpecFile string) (*codedeploy.CreateDeploymentOutput, error) {
	deployOut, err := deploy.Client.CreateDeployment(&codedeploy.CreateDeploymentInput{
		ApplicationName:     aws.String(deploy.DeployShift.ApplicationName),
		DeploymentGroupName: aws.String(deploymentGroup),
		Revision: &codedeploy.RevisionLocation{
			RevisionType: aws.String("AppSpecContent"),
			AppSpecContent: &codedeploy.AppSpecContent{
				Content: aws.String(appSpecFile),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating deployment: %w", err)
	}
	return deployOut, nil
}

func (deploy *Deploy) connectCodeDeploy() (*codedeploy.CodeDeploy, error) {
	sess := session.Must(session.NewSession(deploy.Config))
	svc := codedeploy.New(sess)
	return svc, nil
}
