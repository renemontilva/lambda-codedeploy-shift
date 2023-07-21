package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
)

type Pipeline struct {
	Name   string
	Stage  string
	Action codepipeline.ActionDeclaration
	Config *aws.Config
}

func NewPipeline(name string, stage string) (*Pipeline, error) {
	pipeline := Pipeline{
		Name:   name,
		Stage:  stage,
		Config: &aws.Config{},
	}
	err := pipeline.setDeployAction()
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func (p *Pipeline) setDeployAction() error {
	connection, err := p.codepipelineConn()
	if err != nil {
		return fmt.Errorf("Error creating CodePipeline connection: %w", err)
	}
	codepipelineOut, err := connection.GetPipeline(&codepipeline.GetPipelineInput{
		Name: aws.String(p.Name),
	})
	for _, stage := range codepipelineOut.Pipeline.Stages {
		if *stage.Name == p.Stage {
			for _, action := range stage.Actions {
				if *action.Name == "Deploy" {
					p.Action.ActionTypeId = action.ActionTypeId
					p.Action.Configuration = action.Configuration
					p.Action.InputArtifacts = action.InputArtifacts
					p.Action.Name = action.Name
					p.Action.Namespace = action.Namespace
					p.Action.OutputArtifacts = action.OutputArtifacts
					p.Action.Region = action.Region
					p.Action.RoleArn = action.RoleArn
					p.Action.RunOrder = action.RunOrder
				}
			}
		}

	}
	return nil
}

func (p *Pipeline) UpdateDeployAction() error {

	connection, err := p.codepipelineConn()
	if err != nil {
		return fmt.Errorf("Error creating CodePipeline connection: %w", err)
	}
	codepipelineOut, err := connection.GetPipeline(&codepipeline.GetPipelineInput{
		Name: aws.String(p.Name),
	})
	out, _ := connection.GetActionType(&codepipeline.GetActionTypeInput{
		Category: p.Action.ActionTypeId.Category,
		Owner:    p.Action.ActionTypeId.Owner,
		Provider: p.Action.ActionTypeId.Provider,
		Version:  p.Action.ActionTypeId.Version,
	})
	fmt.Println(out, p.Action)
	//actionOutput, err := connection.UpdateActionType(&codepipeline.UpdateActionTypeInput{
	//	ActionType: &codepipeline.ActionTypeDeclaration{
	//		Executor: &codepipeline.ExecutorConfiguration{
	//			Configuration: p.Action.Configuration,
	//			Type:          p.Action.ActionTypeId.Category,
	//		},
	//	},
	//})
	//if err != nil {
	//	return fmt.Errorf("Error updating CodePipeline deploy action: %w", err)
	//}
	//fmt.Println(actionOutput)
	return nil
}

// returns a connection to CodePipeline service.
func (p *Pipeline) codepipelineConn() (*codepipeline.CodePipeline, error) {
	sess, err := session.NewSession(p.Config)
	if err != nil {
		return nil, fmt.Errorf("Error creating session: %w", err)
	}
	return codepipeline.New(sess), nil
}
