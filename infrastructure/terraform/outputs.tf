output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "region" {
  description = "AWS region"
  value       = var.region
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = module.eks.cluster_name
}


output "ecr_name" {
  description = "Docker images registry name"
  value       = aws_ecr_repository.polls_images.name
}

output "ecr_url" {
  description = "Docker images registry url"
  value       = aws_ecr_repository.polls_images.repository_url
}
