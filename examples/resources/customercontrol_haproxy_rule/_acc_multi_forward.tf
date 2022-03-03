resource "customercontrol_haproxy_rule" "multi-forward" {
  domain_name = "terraform-provider-test.amcsgroup.io"
  setup_kind  = "multi-forward"

  setup_configuration_multi_forward {
    set_host     = false
    load_balance = "Source"

    servers {
      url    = "grafana-dev.amcsgroup.io"
      is_ssl = true
      port   = 443
    }
    servers {
      url    = "grafana.amcsgroup.io"
      is_ssl = true
      port   = 443
    }
  }
}
