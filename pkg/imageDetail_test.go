package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageDetail(t *testing.T) {
	imageDetail, err := NewImageDetail("testdata/tmp", "imageDetail.json")

	if err != nil {
		t.Errorf("Error returned from NewImageDetail: %s", err)
	}
	assert.Equal(t, "12345668.dkr.ecr.us-east-1.amazonaws.com/ecsdeploy:latest", imageDetail.ImageURI)
}
