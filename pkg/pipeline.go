package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
)

type Pipeline struct {
	JobDetails codepipeline.JobDetails
	*codepipeline.CodePipeline
	Config *aws.Config
}

type optionPipeline func(*Pipeline) error

func NewPipeline(jobID string, options ...optionPipeline) (*Pipeline, error) {
	pipeline := &Pipeline{
		Config: &aws.Config{},
		JobDetails: codepipeline.JobDetails{
			Id: aws.String(jobID),
		},
	}

	for _, option := range options {
		err := option(pipeline)
		if err != nil {
			return nil, fmt.Errorf("error creating pipeline: %w", err)
		}
	}

	pipeline.connectCodePipeline()
	return pipeline, nil
}

func (p *Pipeline) connectCodePipeline() {
	sess := session.Must(session.NewSession(p.Config))
	p.CodePipeline = codepipeline.New(sess)
}

func (p *Pipeline) JobSucceeded() {
	p.CodePipeline.PutJobSuccessResult(&codepipeline.PutJobSuccessResultInput{
		JobId: p.JobDetails.Id,
	})
}

func (p *Pipeline) JobFailed(message string) {
	p.CodePipeline.PutJobFailureResult(&codepipeline.PutJobFailureResultInput{
		JobId: p.JobDetails.Id,
		FailureDetails: &codepipeline.FailureDetails{
			Message: aws.String(message),
			Type:    aws.String("JobFailed"),
		},
	})
}
