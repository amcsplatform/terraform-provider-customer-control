# Terraform Provider CustomerControl

- Provider documentation: [https://registry.terraform.io/providers/amcsplatform/customercontrol/latest/docs](https://registry.terraform.io/providers/amcsplatform/customercontrol/latest/docs)

## Requirements 
- Terraform 0.12+

## Development
If you're new to provider development, a good place to start is the [Extending Terraform](https://www.terraform.io/docs/extend/index.html) docs.

Set up your local environment by installing [Go](https://golang.org/). 

### Building
```shell
make build
```

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