
# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
#   AZDO_GITHUB_SERVICE_CONNECTION_PAT
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_github" "github_serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "GitHub Service Connection"
  auth_oauth {
    oauth_configuration_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "azuredevops_build_definition" "nightly_build" {
  project_id      = azuredevops_project.project.id
  agent_pool_name = "Azure Pipelines"
  name            = "Nightly Build"
  path            = "\\"

  repository {
    repo_type             = "GitHub"
    repo_id               = "microsoft/terraform-provider-azuredevops"
    branch_name           = "master"
    yml_path              = ".azdo/azure-pipeline-nightly.yml"
    service_connection_id = azuredevops_serviceendpoint_github.github_serviceendpoint.id
  }
}

# Example Service Hook: Webhook for Git push events
resource "azuredevops_servicehook_subscription" "git_push_webhook" {
  project_id         = azuredevops_project.project.id
  publisher_id       = "tfs"
  event_type         = "git.push"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    # Filter for pushes to master branch only
    branch = "refs/heads/master"
  }

  consumer_inputs = {
    # Replace with your webhook URL
    url = "https://webhook.example.com/git-push"
  }

  resource_version = "1.0"
  status          = "enabled"
}

# Example Service Hook: Webhook for build completion
resource "azuredevops_servicehook_subscription" "build_complete_webhook" {
  project_id         = azuredevops_project.project.id
  publisher_id       = "tfs"
  event_type         = "build.complete"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    # Filter for builds from our build definition
    buildDefinition = azuredevops_build_definition.nightly_build.id
  }

  consumer_inputs = {
    # Replace with your webhook URL
    url = "https://webhook.example.com/build-complete"
  }

  resource_version = "1.0"
  status          = "enabled"
}
