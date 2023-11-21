output "region" {
  description = "AWS region"
  value       = var.region
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = module.eks.cluster_name
}

output "api_ecr_url" {
  description = "API - Docker image registry url"
  value       = aws_ecr_repository.polls_api.repository_url
}

output "front_ecr_url" {
  description = "Frontend - Docker image registry url"
  value       = aws_ecr_repository.polls_front.repository_url
}

output "jwt_key" {
  description = "JWT Key used to sign auth tokens in app"
  value       = random_password.jwt_key.result
  sensitive   = true
}

output "db_password" {
  description = "Password used in production database"
  value       = random_password.db.result
  sensitive   = true
}
