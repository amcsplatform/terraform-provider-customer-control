terraform {
  required_providers {
    customercontrol = {
      version = "0.0.1"
      source  = "amcsgroup.com/amcs/customercontrol"
    }
  }
}

provider "customercontrol" {
  url = "https://customercontrol-dev.amcsgroup.io"
  private_key = ""
}

data "customercontrol_haproxy_domain" "test" {
  domain_name = "d1-p83-svc-publisher-proxy.amcsplatform.com"
}

output "domain_id" {
  value = data.customercontrol_haproxy_domain.test.id
}

output "domain_name" {
  value = data.customercontrol_haproxy_domain.test.domain_name
}

output "domain_valid_until" {
  value = data.customercontrol_haproxy_domain.test.valid_until
}

