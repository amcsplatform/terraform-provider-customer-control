terraform {
  required_providers {
    customercontrol = {
      version = "0.1.21"
      source  = "amcsplatform/customercontrol"
    }
  }
}

provider "customercontrol" {
  url         = ""
  private_key = ""
}

resource "customercontrol_haproxy_rule" "simple-forward" {
  domain_name = "terraform-provider-test.amcsgroup.io"
  setup_kind  = "simple-forward"

  setup_configuration {
    backend      = "grafana-dev.amcsgroup.io"
    is_ssl       = true
    backend_port = 443
    set_host     = true
  }
}

//resource "customercontrol_haproxy_rule" "multi-forward" {
//  domain_name = "test.example.com"
//  setup_kind  = "multi-forward"
//
//  setup_configuration_multi_forward {
//    server {
//      url    = "test.example.io"
//      is_ssl = true
//      port   = 443
//    }
//    server {
//      url    = "test2.example.io"
//      is_ssl = true
//      port   = 443
//    }
//  }
//}
