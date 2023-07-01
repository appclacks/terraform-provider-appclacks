resource "appclacks_healthcheck_http" "test_http" {
  name = "test-http"
  interval = "30s"
  timeout = "5s"
  description = "example http health check"
  labels = {
    "env": "prod"
  }
  query = {
    "param1": "value"
  }
  target = "api.appclacks.com"
  port = 443
  protocol = "https"
  method = "GET"
  path = "/healthz"
  valid_status = [200]
  enabled = true
  host = "api.appclacks.com"
}
