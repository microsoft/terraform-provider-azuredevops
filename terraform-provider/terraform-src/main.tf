
provider "azuredevops" {
  version = ">= 0.0.1"
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
