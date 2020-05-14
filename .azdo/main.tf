
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
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "GitHub Service Connection"

  auth_personal {
    # personalAccessToken = "..." Or set with `AZDO_GITHUB_SERVICE_CONNECTION_PAT` env var
  }
}

resource "azuredevops_serviceendpoint_dockerhub" "dockerhub_serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "DockerHub Service Connection"

  # docker_username = "..." - Or set with `AZDO_DOCKERHUB_SERVICE_CONNECTION_USERNAME` env var
  # docker_email    = "..." - Or set with `AZDO_DOCKERHUB_SERVICE_CONNECTION_EMAIL` env var
  # docker_password = "..." - Or set with `AZDO_DOCKERHUB_SERVICE_CONNECTION_PASSWORD` env var
}

resource "azuredevops_build_definition" "nightly_build" {
  project_id      = azuredevops_project.project.id
  agent_pool_name = "Hosted Ubuntu 1604"
  name            = "Nightly Build"

  repository {
    repo_type             = "GitHub"
    repo_id               = "microsoft/terraform-provider-azuredevops"
    repo_name             = "microsoft/terraform-provider-azuredevops"
    branch_name           = "master"
    yml_path              = ".azdo/azure-pipeline-nightly.yml"
    service_connection_id = azuredevops_serviceendpoint_github.github_serviceendpoint.id
  }
}
