package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type TaskDef ecs.TaskDefinition

type optionTaskDef func(*TaskDef) error

func NewTaskDef(folderPath, fileName string, options ...optionTaskDef) (TaskDef, error) {
	var taskDef TaskDef
	err := convertJSONToType(filepath.Join(folderPath, fileName), &taskDef)
	if err != nil {
		return taskDef, fmt.Errorf("error returned from NewTaskDef: %w", err)
	}

	for _, option := range options {
		err := option(&taskDef)
		if err != nil {
			return TaskDef{}, err
		}
	}
	return taskDef, nil
}

func WithImageURI(imageURI string) optionTaskDef {
	return func(td *TaskDef) error {
		td.ContainerDefinitions[0].Image = &imageURI
		return nil
	}
}

func (td *TaskDef) NewRevision(config *aws.Config) error {
	ecsClient := td.connectECS(config)
	taskDefRevision := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: td.ContainerDefinitions,
		Cpu:                  td.Cpu,
		ExecutionRoleArn:     td.ExecutionRoleArn,
		Family:               td.Family,
		Memory:               td.Memory,
		NetworkMode:          td.NetworkMode,
		PlacementConstraints: td.PlacementConstraints,
		TaskRoleArn:          td.TaskRoleArn,
		Volumes:              td.Volumes,
	}
	newRevision, err := ecsClient.RegisterTaskDefinition(
		taskDefRevision,
	)
	if err != nil {
		return fmt.Errorf("error registering task definition: %w", err)
	}
	td.TaskDefinitionArn = newRevision.TaskDefinition.TaskDefinitionArn
	return nil
}

func (td *TaskDef) connectECS(config *aws.Config) *ecs.ECS {
	sess := session.Must(session.NewSession(config))
	ecs := ecs.New(sess)
	return ecs
}
