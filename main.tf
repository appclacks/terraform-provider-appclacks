terraform {
  required_providers {
    appclacks = {
      version = "~> 0.0.1"
      source = "terraform-appclacks.com/appclacksprovider/appclacks"
    }
  }
}

provider "appclacks" {

}

resource "appclacks_healthcheck_dns" "test_dns" {
  name = "test-tf"
  interval = "30s"
  timeout = "3s"
  description = "hello appclacks terraform provider !"
  labels = {
    "env": "prod"
  }
  domain = "mcorbin.fr"
}
