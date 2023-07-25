# Lambda Function with Shift Deployment Strategy


This repository contains the source code for a Lambda function that allows you to shift deployment strategies based on a descriptive .yaml file. The Lambda function uses the information provided in the deployShift.yaml file to modify its deployment behavior, allowing for more flexibility and control over the deployment process.

## Setup and Deployment

1. Clone the repository.

2. Deploy the Lambda function using the provided terraform code

3. Add and Configure the deployShift.yaml file with the desired deployment strategy options in your application's source code



## deployShift.yaml File Structure

To configure the deployment strategies for AWS CodeDeploy Blue/Green deployments, create or modify the deployShift.yaml file with the following structure:

### Blue/Green Deployment
```yaml
applicationName: CodeployApp
deployStrategies:
  - name: Normal
    deploymentGroup: CanaryNormal
  - name: Fast
    deploymentGroup: CanaryFast
  - name: Slow
    deploymentGroup: CanarySlow
preferredStrategy: Normal
```

- **applicationName**: Replace CodeDeployApp with the name of your AWS CodeDeploy application. This is the application to which you want to deploy using the defined strategies. 

- **deployStrategies**: This section defines the deployment strategies and their corresponding deployment groups. You can customize the deployment strategies by adding or removing entries as needed. Each entry contains the following attributes:

    * **name**: The name of the deployment strategy. Choose a descriptive name that represents the intended behavior of the strategy.
    * **deploymentGroup**: The name of the AWS CodeDeploy deployment group associated with the strategy. Make sure you have created these deployment groups in AWS CodeDeploy. 

- **preferredStrategy**: This attribute specifies the preferred deployment strategy to use. The value should match one of the name attributes from the deployStrategies section. The preferred strategy will be used by default if no strategy is explicitly specified during the deployment.



## Usage
To use the YAML file for AWS CodeDeploy Blue/Green deployments, follow these steps:

Save the YAML content in a deployShift.yaml file, in the root of your source code repository.

Customize the YAML content based on your specific application and deployment group names. Ensure that you have already set up the required AWS CodeDeploy application and deployment groups.

Create a deploy stage on codepipeline and add an invoke action to include the lambda function.


## Local Development

To work on the project locally, follow the steps config section:

* [Configuration](config/README.md)
