resource "aws_route53_zone" "main" {
  name = "choigonyok.com"
}

resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "www.choigonyok.com"
  type    = "A"
  alias {
    name = aws_lb.blog.dns_name
    zone_id = aws_lb.blog.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ci" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "ci.choigonyok.com"
  type    = "A"
  alias {
    name = aws_lb.blog.dns_name
    zone_id = aws_lb.blog.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "cd" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "cd.choigonyok.com"
  type    = "A"
  alias {
    name = aws_lb.blog.dns_name
    zone_id = aws_lb.blog.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53domains_registered_domain" "choigonyok_com" {
  domain_name = "choigonyok.com"

  name_server {
    name = "${aws_route53_zone.main.name_servers.0}"
  }

  name_server {
    name = "${aws_route53_zone.main.name_servers.1}"
  }

  name_server {
    name = "${aws_route53_zone.main.name_servers.2}"
  }

  name_server {
    name = "${aws_route53_zone.main.name_servers.3}"
  }
}