# azuredevops_build_definition
Manages a Build Definition within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_azure_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "Sample Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  repository {
    repo_type   = "TfsGit"
    repo_name   = azuredevops_azure_git_repository.repository.name
    branch_name = azuredevops_azure_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  # Until https://github.com/microsoft/terraform-provider-azuredevops/issues/170, these are assumed
  # to already exist in the project.
  variables_groups = [1, 2, 3]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `name` - (Optional) The name of the build definition.
* `agent_pool_name` - (Optional) The agent pool that should execute the build. Defaults to `Hosted Ubuntu 1604`.
* `repository` - (Required) A `repository` block as documented below.
* `variable_groups` - (Optional) A list of variable group IDs (integers) to link to the build definition.

`repository` block supports the following:

* `branch_name` - (Optional) The branch name for which builds are triggered. Defaults to `master`.
* `repo_name` - (Required) The name of the repository.
* `repo_type` - (Optional) The repository type. Valid values: `GitHub` or `TfsGit`. Defaults to `Github`.
* `service_connection_id` - (Optional) The service connection ID. Used if the `repo_type` is `GitHub`.
* `yml_path` - (Required) The path of the Yaml file describing the build definition.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the build definition
* `revision` - The revision of the build definition

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Build Definitions](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/definitions?view=azure-devops-rest-5.1)

## Import

Not supported