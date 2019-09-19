
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
