resource "aws_vpc" "lab" { // Create a VPC with DNS support enabled
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.project_name}-vpc"
  }
}

resource "aws_s3_bucket" "lab_bucket" {
  # bucket_prefix ensures a unique name by adding a random string to the end
  bucket_prefix = "${var.project_name}-"
  
  tags = {
    Name        = "lab-franklin-bucket-1"
    Environment = "lab"
  }
}

resource "aws_s3_bucket_public_access_block" "lab_bucket_access" {
  bucket                  = aws_s3_bucket.lab_bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}