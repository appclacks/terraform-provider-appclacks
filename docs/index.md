---
page_title: "Provider: Appclacks"
description: |-
  The Appclacks provider allows the use of the Appclacks cloud monitoring platform within Terraform configurations.
---

# Appclacks Provider

The "Appclacks" provider allows the use of the Appclacks cloud monitoring platform within Terraform configurations.

Appclacks is a cloud monitoring platform allowing users to create, manage and run various health checks to monitor the health of websites and infrastructures.

Appclacks can monitor public endpoints by executing configured health checks on them from multiple point of presences. But unlike alternatives,  Appclacks can also be used to monitor private infrastructures.

The prober used by Appclacks named Cabourotte is a free software that you can host on your private infrastructure and plug on the Appclacks API to autoconfigure it.

## Example Usage

```terraform
terraform {
  required_providers {
    appclacks = {
      source = "appclacks/appclacks"
    }
  }
  required_version = ">= 0.1.0"
}

// You can also export the APPCLACKS_ORGANIZATION_ID and
// the APPCLACKS_TOKEN environment variables to configure
// authentication or configure the local Appclacks
// configuration file
provider "appclacks" {
  organization_id = ""
  token = ""
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
3. Local configuration file

To know more about authentication and about how generating an Appclacks API token, please refer to the [Appclacks documentation](https://www.doc.appclacks.com/getting-started/#authentication).

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform
configuration and risks secret leakage should this file ever be committed to a
public version control system.

Credentials can be provided by adding an `organization_id`, `token`, and optionally `api_url`, to the `appclacks` provider block.

Usage:

```terraform
provider "appclacks" {
  organization_id = "my-organization_id"
  token = "my-token"
}
```

### Environment Variables

Credentials can also be provided by setting the `APPCLACKS_ORGANIZATION_ID`, `APPCLACKS_TOKEN`, and optionally `APPCLACKS_API_URL` environment variables.

For example:

```terraform
provider "appclacks" {}
```

```bash
export APPCLACKS_ORGANIZATION_ID="my-organization_id"
export APPCLACKS_TOKEN="my-token"
terraform plan
```

## Appclacks Configuration Reference

| Name             | Description               |  Type  |            Default            | Required |
| ---------------- | ------------------------- | :----: | :---------------------------: | :------: |
| api\_url         | Appclacks API URL         | string | `"https://api.appclacks.com"` |    no    |
| organization\_id | Appclacks organization ID | string |              n/a              |   yes    |
| token            | Appclacks API token       | string |              n/a              |   yes    |
