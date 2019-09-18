
# Make sure to set AZDO_PERSONAL_ACCESS_TOKEN to configure the provider correctly.
provider "azuredevops" {
  version = ">= 0.0.1"

  # This can optionally be provided using the AZDO_ORG_SERVICE_URL env var
  org_service_url = "https://dev.azure.com/niiodice"
}

resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description" # (OPTIONAL DEFAULT "")
  visibility         = "private"                  # public, private (OPTIONAL DEFAULT: private)
  version_control    = "Git"                      # Git, Tfvc (OPTIONAL DEFAULT: git)
  work_item_template = "Scrum"                    # Scrum, Agile, Basic, CMMI (OPTIONAL DEFAULT: Agile)

  #TODO support Custom templates (process templates)
}

resource "azuredevops_pipeline" "pipeline" {
  project_id    = azuredevops_project.project.project_id
  pipeline_name = "Test Pipeline"

  repository {
    repo_type             = "GitHub"
    repo_name             = "nmiodice/terraform-azure-devops-hack"
    branch_name           = "master"
    yml_path              = "azdo-api-samples/azure-pipeline.yml"
    service_connection_id = "..."
  }
}
