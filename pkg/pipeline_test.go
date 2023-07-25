package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPipeline(t *testing.T) {
	t.Run("should return a Pipeline struct", func(t *testing.T) {
		pipeline, err := NewPipeline("1234567890")
		assert.Nil(t, err)
		assert.NotNil(t, pipeline)
		assert.Equal(t, "1234567890", *pipeline.JobDetails.Id)
	})

}
