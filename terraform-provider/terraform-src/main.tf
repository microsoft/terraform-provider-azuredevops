
provider "azuredevops" {
  version = ">= 0.0.1"
  # provide via env var as AZDO_PERSONAL_ACCESS_TOKEN=<my personal access token>
  personal_access_token = "foo"
  # provide via env var as AZDO_ORG_SERVICE_URL=<my org's service url>
  org_service_url = "https://dev.azure.com/chzipp"
}

resource "azuredevops_foo" "nicks_example_resource" {
  fookey = "fooValue"
}

# resource "azuredevops_$RESOURCE" "$LOGICAL_NAME" {
#   <resource configuration goes here>
# }

# resource "azuredevops_project" "nicks_project" {
#   organization = "..."
#   project_name = "..."
#   ...
# }

# resource "azuredevops_pipeline" "nicks_pipeline" {
#   organization = "..."
#   project_id = azuredevops_project.nicks_project.id
#   ...
# }
