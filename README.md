# Terraform Provider for Azure DevOps

[![Build Status](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_apis/build/status/Nightly%20Build?branchName=master)](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_build/latest?definitionId=27&branchName=master)
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

resource "azuredevops_build_definition" "build_definition" {
  project_id      = azuredevops_project.project.id
  name            = "Test Pipeline"
  agent_pool_name = "Hosted Ubuntu 1604"

  repository {
    repo_type             = "GitHub"
    repo_name             = "nmiodice/terraform-azure-devops-hack"
    branch_name           = "master"
    yml_path              = "azdo-api-samples/azure-pipeline.yml"
    service_connection_id = "1a0e1da9-57a6-4470-8e96-160a622c4a17" # Note: Eventually this will come from a GitHub Service Connection resource...  
  }
}
```

# Contributing

Interested in contributing to the provider? Great, we need your help. Get started by reading the [contributing](./docs/contributing.md) document.
