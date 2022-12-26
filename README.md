# Appclacks Terraform provider

Official Terraform provider for [Appclacks](https://appclacks.com/).

See the [documentation](https://registry.terraform.io/providers/appclacks/appclacks/latest/docs) to learn how to use it.

## Launch the tests

Tests will create real resources on your account.

```
export APPCLACKS_ORGANIZATION_ID="<org-id>"
export APPCLACKS_TOKEN='<token>'
export TF_ACC=true
go test -v -race ./...
```
