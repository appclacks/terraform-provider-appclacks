# Appclacks Terraform provider

Official Terraform provider for [Appclacks](https://appclacks.com/).

See the [documentation](https://registry.terraform.io/providers/appclacks/appclacks/latest/docs) to learn how to use it.

## Authentication

You need an API token and to set the `APPCLACKS_ORGANIZATION_ID` and `APPCLACKS_TOKEN` environment variables in order to use the Terraform provider.

See the documentation [Authentication](https://www.doc.appclacks.com/getting-started/index.html#authentication) section for more information about API tokens.

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
