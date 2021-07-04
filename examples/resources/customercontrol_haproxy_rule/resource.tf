resource "customercontrol_haproxy_rule" "example" {
  domain_name = "test.example.com"
  setup_kind  = "simple-forward"

  setup_configuration {
    backend      = "redirect.example.io"
    is_ssl       = true
    backend_port = 443
    set_host     = true
  }
}