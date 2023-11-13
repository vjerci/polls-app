
provider "aws" {
  region = var.region

  default_tags {
    tags = {
      environment = "${var.environment}"
      service     = "polls-app"
    }
  }
}

locals {
  name = "polls-${var.environment}"
}

data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"

  name = local.name

  cidr = "10.0.0.0/16"
  azs  = slice(data.aws_availability_zones.available.names, 0, 2)

  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.3.0/24", "10.0.4.0/24"]

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
}

resource "aws_ecr_repository" "polls_images" {
  name                 = local.name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.15.3"

  cluster_name    = local.name
  cluster_version = "1.27"

  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = module.vpc.private_subnets
  cluster_endpoint_public_access = true

  eks_managed_node_group_defaults = {
    ami_type = "BOTTLEROCKET_x86_64"
  }

  eks_managed_node_groups = {
    app = {
      name = "node-group-app"

      instance_types = ["t3.nano"]
      tags = {
        node_group = "app"
      }

      min_size     = 1
      max_size     = 3
      desired_size = 1
    }

    db = {
      name = "node-group-db"
      tags = {
        node_group = "db"
      }

      instance_types = ["t3.nano"]

      min_size     = 1
      max_size     = 1
      desired_size = 1
    }
  }
}
