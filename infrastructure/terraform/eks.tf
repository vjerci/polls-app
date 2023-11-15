data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
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

      instance_types = ["t3.small"]

      labels = {
        polls-app = "app"
      }

      tags = {
        node_group = "app"
      }

      min_size     = 1
      max_size     = 2
      desired_size = 1
    }

    db = {
      name = "node-group-db"
      labels = {
        polls-app = "db"
      }
      tags = {
        node_group = "db"
      }

      instance_types = ["t3.small"]

      min_size     = 1
      max_size     = 1
      desired_size = 1
    }
  }
}
