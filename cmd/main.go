package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	handlers "github.com/renemontilva/lambda-codedeploy-shift/pkg"
)

func main() {
	lambda.Start(handlers.Handler)
}
