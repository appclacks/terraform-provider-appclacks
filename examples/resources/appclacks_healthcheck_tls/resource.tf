resource "appclacks_healthcheck_tls" "test_tls" {
  name = "test-tls"
  interval = "30s"
  timeout = "5s"
  description = "example tls health check"
  labels = {
    "env": "prod"
  }
  target = "appclacks.com"
  port = 443
  expiration_delay = "168h"
  server_name = "appclacks.com"
  enabled = true
}
