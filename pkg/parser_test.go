package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertYAMLToAppSpecType(t *testing.T) {
	testCases := []struct {
		desc     string
		file     string
		dataType AppSpec
		want     string
	}{
		{
			desc:     "Test appspec.yaml structure",
			file:     "testdata/tmp/appspec.yaml",
			dataType: AppSpec{},
			want:     "<TASK_DEFINITION>",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := convertYAMLToType(tC.file, &tC.dataType)
			assert.Nil(t, err)
			assert.NotNil(t, tC.dataType)
			assert.Equal(t, tC.want, tC.dataType.Resources[0].TargetService.Properties.TaskDefinition)
		})
	}
}

func TestConvertJSONToTaskDefType(t *testing.T) {
	testCases := []struct {
		desc     string
		file     string
		dataType TaskDef
		want     string
	}{
		{
			desc:     "Test taskdef_dev.json structure",
			file:     "testdata/tmp/taskdef_dev.json",
			dataType: TaskDef{},
			want:     "<IMAGE1_NAME>",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := convertJSONToType(tC.file, &tC.dataType)
			assert.Nil(t, err)
			assert.NotNil(t, tC.dataType)
			assert.Equal(t, tC.want, *tC.dataType.ContainerDefinitions[0].Image)
		})
	}
}
