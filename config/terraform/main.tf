
# S3 
resource "aws_s3_bucket" "codepipeline_bucket" {
    bucket_prefix = "lambda-codedeploy-"
}

resource "aws_s3_bucket_acl" "codepipeline_bucket_acl" {
  bucket = aws_s3_bucket.codepipeline_bucket.id
  acl    = "private"
}

resource "aws_codestarconnections_connection" "codepipeline_connection" {
  name          = "lambda-codedeploy-connection"
  provider_type = "GitHub"
}

# Codepipeline
