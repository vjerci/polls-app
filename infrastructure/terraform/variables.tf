variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}


# this is closly tied to region
# a list can be found here: https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html
variable "aws_images" {
  description = "AWS eks addons images repository"
  type        = string
  default     = "602401143452.dkr.ecr.us-east-1.amazonaws.com"
}


variable "environment" {
  description = "enviornment ex. development, staging, production"
  type        = string
  default     = "development"
}



locals {
  name = "polls-${var.environment}"
  app  = "polls"
}
