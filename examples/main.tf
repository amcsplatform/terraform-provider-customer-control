terraform {
  required_providers {
    customercontrol = {
      version = "0.0.14"
      source  = "amcsgroup.com/amcs/customercontrol"
    }
  }
}

provider "customercontrol" {
  url         = ""
  private_key = ""
}

resource "customercontrol_haproxy_rule" "test" {
  domain_name = "test.example.com"
  setup_kind  = "simple-forward"

  setup_configuration {
    backend      = "test.example.io"
    is_ssl       = true
    backend_port = 443
    set_host     = true
  }
}
