# Appclacks Terraform provider

Official Terraform provider for [https://appclacks.com/](https://appclacks.com/).

## Launch the tests

Tests will create real resources on your account.

```
export APPCLACKS_ORGANIZATION_ID="<org-id>"
export APPCLACKS_TOKEN='<token>'
export TF_ACC=true
go test -v -race ./...
```
