# Appclacks Terraform provider

Official Terraform provider for [Appclacks](https://appclacks.com/).

The [Terraform provider documentation](https://registry.terraform.io/providers/appclacks/appclacks/latest/docs) explains how to use it.

## Examples

The `examples` directory provides configuration examples.

## Launch the tests

Tests will create real resources on your account.

```
export APPCLACKS_ORGANIZATION_ID="<org-id>"
export APPCLACKS_TOKEN='<token>'
export TF_ACC=true
go test -v -race ./...
```

## Generate documentation

```
go generate ./...
```
