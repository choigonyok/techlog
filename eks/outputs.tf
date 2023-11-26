output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_security_group_id" {
  description = "Security group ids attached to the cluster control plane"
  value       = module.eks.cluster_security_group_id
}

output "region" {
  description = "AWS region"
  value       = var.region
}

output "cluster_name" {
  description = "Kubernetes Cluster Name"
  value       = module.eks.cluster_name
}  

output "vpc_id" {
  description = "Kubernetes Cluster VPC id"
  value       = module.vpc.default_vpc_id
}  

output "efs_id" {
  value = aws_efs_file_system.jenkins_volume.id
}

output "cluster_oidc_issuer_url" {
  value = module.eks.cluster_oidc_issuer_url
}

output "cluster_oidc_provider_arn" {
  value = module.eks.oidc_provider_arn
}

output "cluster_oidc_provider" {
  value = module.eks.oidc_provider
}

output "load_balancer_dns" {
  value = aws_lb.blog.dns_name
}

output "target_group_http_arn" {
  value = aws_lb_target_group.http.arn
}

output "target_group_https_arn" {
  value = aws_lb_target_group.https.arn
}

output "host_zone_name_servers" {
  value = aws_route53_zone.main.name_servers
}

output "test1" {
  value = aws_lb.blog.ip_address_type
}

output "test2" {
  value = aws_lb.blog.dns_name
}