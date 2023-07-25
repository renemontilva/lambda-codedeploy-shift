package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type AppSpec struct {
	Version   string     `yaml:"version"`
	Resources []Resource `yaml:"Resources"`
}

type Resource struct {
	TargetService TargetService `yaml:"TargetService"`
}

type TargetService struct {
	Type       string     `yaml:"Type"`
	Properties Properties `yaml:"Properties"`
}

type Properties struct {
	TaskDefinition   string           `yaml:"TaskDefinition"`
	LoadBalancerInfo LoadBalancerInfo `yaml:"LoadBalancerInfo"`
}

type LoadBalancerInfo struct {
	ContainerName string `yaml:"ContainerName"`
	ContainerPort int    `yaml:"ContainerPort"`
}

type optionAppSpec func(*AppSpec) error

func NewAppSpec(folderPath, fileName string, options ...optionAppSpec) (AppSpec, error) {
	appSpec := AppSpec{}
	yamlFile, err := os.ReadFile(filepath.Join(folderPath, fileName))
	if err != nil {
		return appSpec, fmt.Errorf("error reading YAML file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &appSpec)
	if err != nil {
		return appSpec, fmt.Errorf("error unmarshalling YAML file: %w", err)
	}

	for _, option := range options {
		err := option(&appSpec)
		if err != nil {
			return AppSpec{}, err
		}
	}
	return appSpec, nil
}

func WithTaskDefinition(taskDefARN string) optionAppSpec {
	return func(app *AppSpec) error {
		app.Resources[0].TargetService.Properties.TaskDefinition = taskDefARN
		return nil
	}
}

func (app *AppSpec) String() (string, error) {
	yamlBytes, err := yaml.Marshal(app)
	if err != nil {
		return fmt.Sprintf("Error marshalling YAML: %v", err), nil
	}
	return string(yamlBytes), nil
}
