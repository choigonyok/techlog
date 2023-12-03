data "aws_acm_certificate" "techlog" {
  domain      = "www.choigonyok.com"
  types       = ["AMAZON_ISSUED"]
}
