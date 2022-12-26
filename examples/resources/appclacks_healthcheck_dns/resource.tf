resource "appclacks_healthcheck_dns" "test_dns" {
  name = "test-dns"
  interval = "30s"
  timeout = "5s"
  description = "example dns health check"
  labels = {
    "env": "prod"
  }
  domain = "appclacks.com"
  enabled = true
}
