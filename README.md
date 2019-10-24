# Terraform Provider for Azure DevOps

[![Gitter](https://badges.gitter.im/terraform-provider-azuredevops/community.svg)](https://gitter.im/terraform-provider-azuredevops/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Build Status](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_apis/build/status/Nightly%20Build?branchName=master)](https://dev.azure.com/terraform-azdo/terraform-provider-azuredevops/_build/latest?definitionId=27&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/microsoft/terraform-provider-azuredevops)](https://goreportcard.com/report/github.com/microsoft/terraform-provider-azuredevops)

The AzDO (Azure DevOps) Provider supports Terraform 0.11.x and later - but Terraform 0.12.x is recommended.

* [Terraform Website](https://www.terraform.io)
* [AzDO Website](https://azure.microsoft.com/en-us/services/devops/)
* [AzDO Provider Usage Examples](./examples/)
* [AzDO Provider Reference](./website/index.md)

Checkout our [Project Roadmap](./docs/roadmap.md).

## Important!
This repository is a work in progress and is not yet suitable for production workloads. Community contributions are welcome.

## Configuration Values

| Environment Variable | Description | Required? | Example |
| --- | --- | --- | --- |
| `AZDO_PERSONAL_ACCESS_TOKEN` | A personal access token that grants access to Azure DevOps APIs within the org specified by `AZDO_ORG_SERVICE_URL` | yes | `d7894a91db7610e39decbe09b2dfd449ed2ed5a` |
| `AZDO_ORG_SERVICE_URL` | URL of the Azure DevOps org in which resources will be provisioned/managed | yes | `https://dev.azure.com/contoso-org` |
| `AZDO_GITHUB_SERVICE_CONNECTION_PAT` | If running the acceptance tests, you will need this defined in order to validate the GitHub Service Connection resource | for acceptance tests only | `a9194a91d75643e39decbe09b2dfd558dd2abca` |
| `AZDO_PRJ_CREATE_DELAY` | Delay (in seconds) to insert after creation of projects. This was determined to be useful based on observed behavior of the AzDO APIs | no | `10` |

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
