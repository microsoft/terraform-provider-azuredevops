
# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
#   AZDO_GITHUB_SERVICE_CONNECTION_PAT
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = "=0.0.998"
    }
  }
}

provider "azuredevops" {
  org_service_url       = var.org_url
  personal_access_token = var.org_token
}

data "azuredevops_group" "tf_padrao_projectadm" {
  project_id = azuredevops_project.project.id
  name       = "Project Administrators"
}

resource "azuredevops_project" "project" {
  name               =  var.project_name
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_project_releases_permissions" "testando" {
  project_id  = azuredevops_project.project.id
  principal   = data.azuredevops_group.tf_padrao_projectadm.descriptor
  permissions = {
    CreateReleases = "Deny"
  }
}
