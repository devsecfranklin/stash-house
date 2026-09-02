variable "project_name" {
  default = "lab-franklin"
}

variable "region" {
  description = "Region for AWS resources"
  type        = string
  default     = "us-west-1"
}

variable "ec2_ssh_key_name" {
  description = "The SSH Key Name"
  type        = string
  default     = "free-tier-ec2-key"
}

variable "ec2_ssh_public_key_path" {
  description = "The local path to the SSH Public Key"
  type        = string
  default     = "~/.ssh/google_compute_engine.pub"
}

variable "vpc_cidr" {
  default = "172.29.0.0/22"
}
