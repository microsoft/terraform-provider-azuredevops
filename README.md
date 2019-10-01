# Terraform Provider for Azure DevOps

[![Build Status](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_apis/build/status/microsoft.terraform-provider-azuredevops?branchName=master)](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_build/latest?definitionId=3&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/microsoft/terraform-provider-azuredevops)](https://goreportcard.com/report/github.com/microsoft/terraform-provider-azuredevops)

The AzDO (Azure DevOps) Provider supports Terraform 0.11.x and later - but Terraform 0.12.x is recommended.

* [Terraform Website](https://www.terraform.io)
* [AzDO Website](https://azure.microsoft.com/en-us/services/devops/)
* [AzDO Provider Usage Examples](./examples/)

## Important!
This repository is a work in progress and is not yet suitable for production workloads. Community contributions are welcome.

## Usage Example

* Installing the provider
```bash
./scripts/build.sh          # build & test provider code
./scripts/local-install.sh  # install provider locally
```

* Using the provider
```hcl
# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_pipeline" "pipeline" {
  project_id    = azuredevops_project.project.id
  pipeline_name = "Test Pipeline"

  repository {
    repo_type             = "GitHub"
    repo_name             = "nmiodice/terraform-azure-devops-hack"
    branch_name           = "master"
    yml_path              = "azdo-api-samples/azure-pipeline.yml"
    service_connection_id = "1a0e1da9-57a6-4470-8e96-160a622c4a17" # Note: Eventually this will come from a GitHub Service Connection resource...  
  }
}
```

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.11.x +
* [Go](https://golang.org/doc/install) version 1.12.x (to build the provider plugin)

If you're on Windows you'll need to install WSL. Other dependencies called out should be installed within WSL:
* [Installing WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10)


## Developing the Provider

* TODO: Fill section out...

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
