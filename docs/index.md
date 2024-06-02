---
page_title: "Provider: Appclacks"
description: |-
  The Appclacks provider allows the use of the Appclacks monitoring platform within Terraform configurations.
---

# Appclacks Provider

The "Appclacks" provider allows the use of the [Appclacks platform](https://appclacks.com) within Terraform configurations.

Appclacks is an open source platform dedicated to observability.
For example, this provider allows  users to create, manage and run various health checks to monitor the health of websites and infrastructures (blackbox monitoring)

The prober used by Appclacks named Cabourotte is a free software that you can host on your private infrastructure and plug on the Appclacks API to autoconfigure it.

More information about Appclacks can be found in the [official documentation](https://www.doc.appclacks.com/).

## Example Usage

```terraform
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
```

## Authentification and Configuration

Configuration for the Appclacks provider can be provided in the following ways:

1. Parameters in the provider configuration
2. Environment variables

For obtaining an Appclacks API token, please refer to the [Appclacks documentation](https://www.doc.appclacks.com/getting-started/#authentication).

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform
configuration and risks secret leakage should this file ever be committed to a
public version control system.

Credentials can be provided by adding an `api_endpoint`, `username`, and `password`, to the `appclacks` provider block.

TLS options can be set using the `tls_key`, `tls_cert`, `tls_cacert` and `insecure` options.

Usage:

```terraform
provider "appclacks" {
  api_endpoint = "https://appclacks.com"
  username = "my-user"
  password = "my-password"
}
```

### Environment Variables

Credentials can also be provided by setting the `APPCLACKS_API_ENDPOINT`, `APPCLACKS_USERNAME`, and optionally `APPCLACKS_PASSWORD` environment variables.

TLS options can be set with `APPCLACKS_TLS_KEY`, `APPCLACKS_TLS_CERT`, `APPCLACKS_TLS_CACERT` and `APPCLACKS_TLS_INSECURE`.

For example:

```terraform
provider "appclacks" {}
```

```bash
export APPCLACKS_API_ENDPOINT="https://appclacks.com"
export APPCLACKS_USERNAME="my-user"
export APPCLACKS_PASSWORD="password"
terraform plan
```