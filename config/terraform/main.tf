# ECR
module "ecr" {
  source  = "terraform-aws-modules/ecr/aws"
  version = "1.6.0"

  repository_name = "lambda-codedeploy-shift"

  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Keep last 30 images",
        selection = {
          tagStatus     = "tagged",
          tagPrefixList = ["v"],
          countType     = "imageCountMoreThan",
          countNumber   = 30
        },
        action = {
          type = "expire"
        }
      }
    ]
  })

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

# S3 
resource "aws_s3_bucket" "codepipeline_bucket" {
  bucket_prefix = "lambda-codedeploy-"
}

resource "aws_codestarconnections_connection" "codepipeline_connection" {
  name          = "lambda-codedeploy-connection"
  provider_type = "GitHub"
}

# Codepipeline

resource "aws_codepipeline" "codepipeline" {
  name     = "lambda-codedeploy-shift-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.codepipeline_bucket.bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        ConnectionArn    = aws_codestarconnections_connection.codepipeline_connection.arn
        FullRepositoryId = "renemontilva/lambda-codedeploy-shift-poc"
        BranchName       = "main"
      }
    }
  }

  stage {
    name = "Build"

    action {
      name            = "Build"
      category        = "Build"
      owner           = "AWS"
      provider        = "CodeBuild"
      version         = "1"
      input_artifacts = ["source_output"]
      output_artifacts = [ "build_output" ]

      configuration = {
        ProjectName = aws_codebuild_project.shift_build.name
      }
    }
  }

  stage {
    name = "Deploy"
    action {
      name            = "Deploy"
      category        = "Invoke"
      owner           = "AWS"
      provider        = "Lambda"
      version         = "1"
      input_artifacts = ["build_output"]

      configuration = {
        FunctionName   = aws_lambda_function.codedeploy_shift_function.function_name
        UserParameters = "dev"
      }
    }
  }
}


# Lambda 
resource "aws_lambda_function" "codedeploy_shift_function" {
  function_name = "lambda-codedeploy-shift"
  timeout       = 30 # seconds
  image_uri     = "${module.ecr.repository_url}:latest"
  package_type  = "Image"

  role = aws_iam_role.codedeploy_shift_function.arn

  environment {
    variables = {
      ENVIRONMENT = "dev"
    }
  }
}


# Codebuild

resource "aws_codebuild_project" "shift_build" {
  name         = "lambda-codedeploy-shift-build"
  description  = "Build for lambda-codedeploy-shift"
  service_role = aws_iam_role.codebuild_role.arn

  artifacts {
    type = "CODEPIPELINE"
  } 

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/standard:7.0"
    type         = "LINUX_CONTAINER"

    environment_variable {
      name  = "ENVIRONMENT"
      value = "dev"
    }
  }

  logs_config {
    s3_logs {
        encryption_disabled = false
        location            = "${aws_s3_bucket.codepipeline_bucket.bucket}/codebuild"
        status = "ENABLED"
    }
  }

  source {
    type            = "CODEPIPELINE"
    buildspec       = "buildspec.yaml"
  }
}
