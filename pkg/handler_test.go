package handlers

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	codePipelineEvent := events.CodePipelineJobEvent{}
	payload, err := os.ReadFile("testdata/codepipeline_event.json")
	if err != nil {
		t.Errorf("Error reading testdata: %s", err)
	}
	err = json.Unmarshal(payload, &codePipelineEvent)
	if err != nil {
		t.Errorf("Error unmarshalling testdata: %s", err)
	}

	err = Handler(context.TODO(), codePipelineEvent)
	if err != nil {
		t.Errorf("Error returned from Handler: %s", err)
	}
}
