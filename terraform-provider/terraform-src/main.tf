
provider "azuredevops" {
  version = ">= 0.0.1"
  # provide via env var as AZDO_PERSONAL_ACCESS_TOKEN=<my personal access token>
  #personal_access_token = "foo"
  # provide via env var as AZDO_ORG_SERVICE_URL=<my org's service url>
  org_service_url = "https://dev.azure.com/niiodice"
}

resource "azuredevops_foo" "nicks_example_resource" {
  fookey = "fooValue"
  project_id = azuredevops_project.nicks_project.project_id
}

resource "azuredevops_pipeline" "pipeline_example" {
  project_id = "..."
  pipeline_name = "sample-pipeline-nick"

  repository {
    repo_type = "GitHub"
    repo_name = "nmiodice/terraform-azure-devops-hack"
    branch_name = "master"
    yml_path = "azdo-api-samples/azure-pipeline.yml"
    service_connection_id = "..."
  }
}

# resource "azuredevops_$RESOURCE" "$LOGICAL_NAME" {
#   <resource configuration goes here>
# }

resource "azuredevops_project" "nicks_project" {
  project_name = "tf_test"
  description = "test project" #(OPTIONAL DEFAULT "")
  visibility = "private" # public, private (OPTIONAL DEFAULT: private)
  version_control = "Git" # Git, Tfvc (OPTIONAL DEFAULT: git)
  work_item_template = "Scrum" # Scrum, Agile, Basic, CMMI (OPTIONAL DEFAULT: Agile)

  #TODO support Custom templates (process templates)
}

# resource "azuredevops_pipeline" "nicks_pipeline" {
#   organization = "..."
#   project_id = azuredevops_project.nicks_project.project_id
#   ...
# }
