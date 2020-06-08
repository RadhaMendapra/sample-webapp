output "app_dns_name" {
 value="http://${aws_lb.lb.dns_name}/builds"
}