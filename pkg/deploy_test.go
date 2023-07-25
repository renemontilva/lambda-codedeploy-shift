package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreferredStrategy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc           string
		deploymetGroup string
		folderPath     string
		deployFile     string
	}{
		{
			desc:           "Test Normal Strategy",
			deploymetGroup: "CanaryNormal",
			folderPath:     "testdata",
			deployFile:     "deployShift.yaml",
		},
		{
			desc:           "Test Fast Strategy",
			deploymetGroup: "CanaryFast",
			folderPath:     "testdata",
			deployFile:     "deployShift_1.yaml",
		},
		{
			desc:           "Test Slow Strategy",
			deploymetGroup: "CanarySlow",
			folderPath:     "testdata",
			deployFile:     "deployShift_2.yaml",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			deployShift, err := NewDeploy(tC.folderPath, DeployWithDeployShiftFileName(tC.deployFile))
			if err != nil {
				t.Errorf("Error returned from NewDeployShift: %s", err)
			}
			deploymentGroup, err := deployShift.PreferredStrategy()
			if err != nil {
				t.Errorf("Error returned from PreferredStrategy: %s", err)
			}
			assert.Equal(t, tC.deploymetGroup, deploymentGroup, "Deployment group should be deployShift")
		})
	}

}
