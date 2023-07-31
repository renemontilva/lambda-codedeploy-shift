package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppSpec(t *testing.T) {
	appSpec, err := NewAppSpec("testdata/tmp", "appspec.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, appSpec)
	assert.Equal(t, "<TASK_DEFINITION>", appSpec.Resources[0].TargetService.Properties.TaskDefinition)

}

func TestWithTaskDefinition(t *testing.T) {
	appSpec, err := NewAppSpec("testdata/tmp", "appspec.yaml",
		WithTaskDefinition("arn:taskdef"),
	)
	if err != nil {
		t.Errorf("Error returned from NewAppSpec %v", err)
	}
	assert.Equal(t, "arn:taskdef", appSpec.Resources[0].TargetService.Properties.TaskDefinition)

}

func TestAppSpecString(t *testing.T) {
	appSpec, err := NewAppSpec("testdata/tmp", "appspec.yaml",
		WithTaskDefinition("arn:taskdef"),
	)
	if err != nil {
		t.Errorf("Error returned from NewAppSpec %v", err)
	}
	appSpecString, _ := appSpec.String()
	assert.Contains(t, appSpecString, "arn:taskdef")

}
