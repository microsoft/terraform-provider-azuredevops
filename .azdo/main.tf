
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name       = "terraform-provider-azuredevops"
  description        = ""
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_github" "github_serviceendpoint" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "GitHub Service Connection"
  authorization {
    scheme = "PersonalAccessToken"
  }
}

resource "azuredevops_serviceendpoint_dockerhub" "dockerhub_serviceendpoint" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "DockerHub Service Connection"
}

resource "azuredevops_build_definition" "nightly_build" {
  project_id      = azuredevops_project.project.id
  agent_pool_name = "Hosted Ubuntu 1604"
  name            = "Nightly Build"

  repository {
    repo_type             = "GitHub"
    repo_name             = "microsoft/terraform-provider-azuredevops"
    branch_name           = "master"
    yml_path              = ".azdo/azure-pipeline-nightly.yml"
    service_connection_id = azuredevops_serviceendpoint_github.github_serviceendpoint.id
  }
}
