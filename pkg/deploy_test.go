package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentStrategy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc           string
		deploymetGroup string
		deployFile     string
	}{
		{
			desc:           "Test Normal Strategy",
			deploymetGroup: "CanaryNormal",
			deployFile:     "testdata/deployShift.yaml",
		},
		{
			desc:           "Test Fast Strategy",
			deploymetGroup: "CanaryFast",
			deployFile:     "testdata/deployShift_1.yaml",
		},
		{
			desc:           "Test Slow Strategy",
			deploymetGroup: "CanarySlow",
			deployFile:     "testdata/deployShift_2.yaml",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			deployShift, err := NewDeployShift(tC.deployFile)
			if err != nil {
				t.Errorf("Error returned from NewDeployShift: %s", err)
			}
			deploymentGroup, err := deployShift.CurrentStrategy()
			if err != nil {
				t.Errorf("Error returned from CurrentStrategy: %s", err)
			}
			assert.Equal(t, tC.deploymetGroup, deploymentGroup, "Deployment group should be deployShift")
		})
	}

}
