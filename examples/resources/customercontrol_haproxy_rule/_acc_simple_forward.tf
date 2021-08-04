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