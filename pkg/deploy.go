package handlers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DeployShift struct {
	ApplicationName  string           `yaml:"applicationName"`
	DeployStrategies []DeployStrategy `yaml:"deployStrategies"`
	ActiveStrategy   string           `yaml:"activeStrategy"`
}

type DeployStrategy struct {
	Name            string `yaml:"name"`
	DeploymentGroup string `yaml:"deploymentGroup"`
}

// NewDeployShift returns a DeployShift struct.
func NewDeployShift(configFile string) (*DeployShift, error) {
	deployShift := DeployShift{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &deployShift)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling YAML file: %w", err)
	}
	return &deployShift, nil
}

func (deploy *DeployShift) findDeploymentGroup() (string, error) {
	for _, deployStrategy := range deploy.DeployStrategies {
		if deployStrategy.Name == deploy.ActiveStrategy {
			return deployStrategy.DeploymentGroup, nil
		}
	}
	return "", fmt.Errorf("Deployment group not found in deployShift.yaml")
}

// CurrentStrategy returns a deploymentGroup name defined in ActiveStrategy field of deployShift.yaml.
func (deploy *DeployShift) CurrentStrategy() (string, error) {
	return deploy.findDeploymentGroup()
}

func (deploy *DeployShift) DeployStrategy() error {
	return nil
}
