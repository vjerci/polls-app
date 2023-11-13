variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}


variable "environment" {
  description = "enviornment ex. development, staging, production"
  type        = string
  default     = "development"
}
