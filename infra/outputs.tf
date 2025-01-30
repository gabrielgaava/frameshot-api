output "lb_dns_name" {
  description = "DNS do Load Balancer"
  value       = aws_lb.frameshot_lb.dns_name
}