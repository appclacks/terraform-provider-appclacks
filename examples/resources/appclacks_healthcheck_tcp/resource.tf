resource "appclacks_healthcheck_tcp" "test_tcp" {
  name = "test-tcp"
  interval = "30s"
  timeout = "5s"
  description = "example tcp health check"
  labels = {
    "env": "prod"
  }
  target = "appclacks.com"
  port = 443
  enabled = true
}
