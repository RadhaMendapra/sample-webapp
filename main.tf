//region and provider info
provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region = "us-west-2"
}

//1 LB
resource "aws_lb" "lb" {
  name = "webapp-lb"
  internal = true
  load_balancer_type = "application"
  security_groups = [aws_security_group.lb_group.id]

  tags = {
    Name = "webapp-lb"
  }
}

//1 EC2 instance
resource "aws_instance" "instance"{
  ami=var.ami_id
  instance_type="t2.micro"
  security_groups = [aws_security_group.instance_group.id]
  key_name = aws_key_pair.key_pair.key_name
}

//security groups
//instance security group
resource "aws_security_group" "instance_group" {
  name = "instance_group"
  description = "App instance security group"

  tags = {
    Name = "instance_group"
  }
}

//security group rule for instance security group
resource "aws_security_group_rule" "app_ingress_http" {
  type = "ingress"
  security_group_id = aws_security_group.instance_group.id

  from_port = 8080
  to_port = 8080
  protocol = "tcp"
  source_security_group_id = aws_security_group.lb_group.id
}

//load balancer security group
resource "aws_security_group" "lb_group" {
  name = "lb_group"
  description = "Load balancer security group"

  tags = {
    Name = "lb_group"
  }
}

//security group rules for lb security group
resource "aws_security_group_rule" "lb_ingress" {
  type = "ingress"
  security_group_id = aws_security_group.lb_group.id

  from_port = 80
  to_port = 80
  protocol = "tcp"
}

resource "aws_security_group_rule" "lb_egress" {
  type = "egress"
  security_group_id = aws_security_group.lb_group.id

  from_port = 8080
  to_port = 8080
  protocol = "tcp"
  source_security_group_id = aws_security_group.instance_group.id
}

resource "aws_security_group" "conn_group" {
  name = "conn_group"
  description = "Server http outbound security group"

  tags = {
    Name = "conn_group"
  }
}

//key pair
resource "tls_private_key" "ssh_key" {
  algorithm = "RSA"
  rsa_bits = 4096
}

resource "aws_key_pair" "key_pair" {
  key_name = "webapp-key"
  public_key = tls_private_key.ssh_key.public_key_openssh
}
