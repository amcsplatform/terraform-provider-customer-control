# Terraform Provider CustomerControl

- Provider documentation: [https://registry.terraform.io/providers/amcsplatform/customercontrol/latest/docs](https://registry.terraform.io/providers/amcsplatform/customercontrol/latest/docs)

## Requirements 
- Terraform 0.12+

## Development
If you're new to provider development, a good place to start is the [Extending Terraform](https://www.terraform.io/docs/extend/index.html) docs.

Set up your local environment by installing [Go](https://golang.org/). 

### Updating documentation
```shell
go generate
```

If the schema is changed, but documentation is not updated, then the pipeline will detect that and fail.

### Building
```shell
make build
```

### Running acceptance tests
```shell
make testacc
```

Tests are defined in `*_test.go` files. They depend on DNS records registered in Azure DNS.
There are currently 2 records required for the tests:
- `terraform-provider-test.amcsgroup.io`, points to `proxy-dev.amcsgroup.io`
- `terraform-provider-test-2.amcsgroup.io`, points to `proxy-dev.amcsgroup.io`

### Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

### Publishing
- Commit to Azure DevOps
- `terraform-provider-customercontrol` ADO pipeline will set build number & Git tag and push code to GitHub
- GitHub Action will trigger on new tag and publish the provider as GitHub release
- Terraform Registry will pick new release up automatically
