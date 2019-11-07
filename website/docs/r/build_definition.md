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
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"

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

## Arugument Reference

The following arguments are supported:

* `agent_pool_name` - (Optional) The agent pool that should execute the build. Defaults to `Hosted Ubuntu 1604`
* `name` - (Optional) The name of the build definition
* `project_id` - (Required) The ID of the project in which to configure the build definition
* `repository` - (Required) A `repository` block a defined below
* `variable_groups` - (Optional) A list of variable group IDs (integers) to link to the build definition

---
A `repository` block supports the following:

* `branch_name` - (Optional) The branch name for which builds are triggered. Defaults to `master`
* `repo_name` - (Required) The name of the repository
* `repo_type` - (Optional) The repository type. Values can be `GitHub` (default) or `TfsGit`.
* `service_connection_id` - (Optional) The service connection ID. Used if the repository type is `GitHub`
* `yml_path` - (Required) The path of the Yaml file describing the build definition


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the build definition
* `revision` - The revision of the build definition

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Build Definitions](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/definitions?view=azure-devops-rest-5.1)

## Import

Not supported