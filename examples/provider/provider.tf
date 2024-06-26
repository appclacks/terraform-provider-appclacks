terraform {
  required_providers {
    appclacks = {
      source = "appclacks/appclacks"
    }
  }
}

// You can also export the APPCLACKS_USERNAME,
// APPCLACKS_PASSWORD and APPCLACKS_API_ENDPOINT variables
// to configure the client.
// See the documentation for more information about authentication: https://www.doc.appclacks.com/getting-started/
provider "appclacks" {
  api_endpoint = ""
  username = ""
  password = ""
}

resource "appclacks_healthcheck_command" "test_command" {
  name = "test-command"
  interval = "30s"
  timeout = "5s"
  description = "example command health check"
  labels = {
    "env": "prod"
  }
  command = "ls"
  arguments = ["/"]
}

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

resource "appclacks_healthcheck_http" "test_http" {
  name = "test-http"
  interval = "30s"
  timeout = "5s"
  description = "example http health check"
  labels = {
    "env": "prod"
  }
  target = "api.appclacks.com"
  port = 443
  protocol = "https"
  method = "GET"
  path = "/healthz"
  valid_status = [200]
  enabled = true
  insecure = true
  server_name = "api.appclacks.com"
}

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
  insecure = false
}

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
