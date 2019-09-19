# Terraform Provider for Azure DevOps

The AzDO (Azure DevOps) Provider supports Terraform 0.11.x and later - but Terraform 0.12.x is recommended.

* [Terraform Website](https://www.terraform.io)
* [AzDO Website](https://azure.microsoft.com/en-us/services/devops/)
* [AzDO Provider Usage Examples](./examples/)

## Usage Example

* Installing the provider
```bash
./build.sh          # build & test provider code
./local-install.sh  # install provider locally
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
  project_id    = azuredevops_project.project.project_id
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
