package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskDef(t *testing.T) {
	taskDef, err := NewTaskDef("testdata/tmp", "taskdef_dev.json")
	if err != nil {
		t.Errorf("Error returned from NewTaskDef: %s", err)
	}
	assert.Equal(t, "<TASK_DEFINITION>", *taskDef.TaskDefinitionArn)
	assert.Equal(t, "<IMAGE1_NAME>", *taskDef.ContainerDefinitions[0].Image)
}

func TestWithImageURI(t *testing.T) {
	taskDef, err := NewTaskDef("testdata/tmp", "taskdef_dev.json",
		WithImageURI("image123"),
	)
	if err != nil {
		t.Errorf("Error returned from NewTaskDef: %s", err)
	}
	assert.Equal(t, "image123", *taskDef.ContainerDefinitions[0].Image)
}
