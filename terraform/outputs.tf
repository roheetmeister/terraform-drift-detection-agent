output "bucket_name" {
  value = aws_s3_bucket.drift_detector.bucket
}

output "bucket_arn" {
  value = aws_s3_bucket.drift_detector.arn
}

output "bucket_region" {
  value = aws_s3_bucket.drift_detector.region
}
