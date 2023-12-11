resource "aws_lb" "blog" {
  name               = "blog"
  internal           = false
  load_balancer_type = "network"
  subnets            = [for subnet in module.vpc.public_subnets : subnet]

  enable_deletion_protection = false

  tags = {
    Environment = "production"
  }
}

resource "aws_lb_target_group" "http" {
  name     = "blog-http-tg"
  port     = 31665
  protocol = "TCP"
  vpc_id   = module.vpc.vpc_id
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.blog.arn
  port              = "80"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.http.arn
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.blog.arn
  port              = "443"
  protocol          = "TLS"
  # certificate_arn   = data.aws_acm_certificate.techlog.arn
  alpn_policy       = "None"
  certificate_arn = aws_acm_certificate_validation.techlog.certificate_arn
  # test
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.http.arn
  }
}