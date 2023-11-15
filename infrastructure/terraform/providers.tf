provider "aws" {
  region = var.region

  default_tags {
    tags = {
      App        = "polls"
      Enviroment = var.environment
    }
  }
}


terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.7.0"
    }
  }

  required_version = "~> 1.3"
}
