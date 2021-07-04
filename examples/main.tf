terraform {
  required_providers {
    customercontrol = {
      version = "0.0.14"
      source  = "amcsgroup.com/amcs/customercontrol"
    }
  }
}

provider "customercontrol" {
  url         = "https://customercontrol-dev.amcsgroup.io"
  private_key = ""
}

resource "customercontrol_haproxy_rule" "test" {
  domain_name = "provider-test.amcsplatform.com"
  setup_kind  = "simple-forward"

  setup_configuration {
    backend      = "grafana.amcsgroup.io"
    is_ssl       = true
    backend_port = 443
    set_host     = true
  }
}
//
//output "domain_id" {
//  value = customercontrol_haproxy_rule.test.id
//}
//
//output "domain_name" {
//  value = customercontrol_haproxy_rule.test.domain_name
//}
//
//output "domain_valid_until" {
//  value = customercontrol_haproxy_rule.test.valid_until
//}

