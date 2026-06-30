terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-south-1"
}

resource "aws_s3_bucket" "drift_detector" {
  bucket = "ai-tf-drift-detector"

  tags = {
    Name        = "ai-tf-drift-detector"
    Project     = "terraform-drift-detection-agent"
    Environment = "development"
    ManagedBy   = "terraform"
  }
}

resource "aws_s3_bucket_versioning" "drift_detector" {
  bucket = aws_s3_bucket.drift_detector.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "drift_detector" {
  bucket = aws_s3_bucket.drift_detector.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "drift_detector" {
  bucket = aws_s3_bucket.drift_detector.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
