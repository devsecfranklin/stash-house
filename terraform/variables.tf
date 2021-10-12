variable "region" {
  description = "AWS Region for deployment, for example \"us-east-1\"."
  type        = string
}

variable "vpc_cidr" {
  description = "Security VPC CIDR"
  type        = string
  default     = "10.100.0.0/16"
}

variable "subnets_cidr" {
  type    = list(string)
  default = ["10,100.0.128/25"]
}

variable "ssh_key_name" {
  description = "name of pub SSH key, manual add to ec2 console"
  type        = string
  default     = "fdiaz-ssh-key"
}
