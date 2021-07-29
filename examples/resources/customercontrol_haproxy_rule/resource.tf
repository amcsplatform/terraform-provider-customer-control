// Simple-forward kind
resource "customercontrol_haproxy_rule" "simple-forward" {
  domain_name = "test.example.com"
  setup_kind  = "simple-forward"

  setup_configuration {
    backend      = "redirect.example.io"
    is_ssl       = true
    backend_port = 443
    set_host     = true
  }
}

// Multi-forward kind
resource "customercontrol_haproxy_rule" "multi-forward" {
  domain_name = "text.example.com"
  setup_kind  = "multi-forward"

  setup_configuration_multi_forward {
    servers {
      url    = "redirect-1.example.io"
      is_ssl = true
      port   = 443
    }
    servers {
      url    = "redirect-2.example.io"
      is_ssl = true
      port   = 443
    }
  }
}
