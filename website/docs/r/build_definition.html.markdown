---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_build_definition"
description: |-
  Manages a Build Definition within Azure DevOps organization.
---

# azuredevops_build_definition

Manages a Build Definition within Azure DevOps.

## Example Usage

### Tfs
```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "Sample Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_variable_group" "vars" {
  project_id   = azuredevops_project.project.id
  name         = "Infrastructure Pipeline Variables"
  description  = "Managed by Terraform"
  allow_access = true

  variable {
    name  = "FOO"
    value = "BAR"
  }
}

resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.repository.id
    branch_name = azuredevops_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  variable_groups = [
    azuredevops_variable_group.vars.id
  ]

  variable {
    name  = "PipelineVariable"
    value = "Go Microsoft!"
  }

  variable {
    name      = "PipelineSecret"
    secret_value     = "ZGV2cw"
    is_secret = true
  }
}
```

### GitHub Enterprise
```hcl
resource "azuredevops_build_definition" "sample_dotnetcore_app_release" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type             = "GitHubEnterprise"
    repo_id               = "<GitHub Org>/<Repo Name>"
    github_enterprise_url = "https://github.company.com"
    branch_name           = "master"
    yml_path              = "azure-pipelines.yml"
    service_connection_id = "..."
  }

}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `name` - (Optional) The name of the build definition.
- `path` - (Optional) The folder path of the build definition.
- `agent_pool_name` - (Optional) The agent pool that should execute the build. Defaults to `Hosted Ubuntu 1604`.
- `repository` - (Required) A `repository` block as documented below.
- `ci_trigger` - (Optional) Continuous Integration Integration trigger.
- `pull_request_trigger` - (Optional) Pull Request Integration Integration trigger.
- `variable_groups` - (Optional) A list of variable group IDs (integers) to link to the build definition.
- `variable` - (Optional) A list of `variable` blocks, as documented below.

`variable` block supports the following:

- `name` - (Required) The name of the variable.
- `value` - (Optional) The value of the variable.
- `secret_value` - (Optional) The secret value of the variable. Used when `is_secret` set to `true`.
- `is_secret` - (Optional) True if the variable is a secret. Defaults to `false`.
- `allow_override` - (Optional) True if the variable can be overridden. Defaults to `true`.

`repository` block supports the following:

- `branch_name` - (Optional) The branch name for which builds are triggered. Defaults to `master`.
- `repo_id` - (Required) The id of the repository. For `TfsGit` repos, this is simply the ID of the repository. For `Github` repos, this will take the form of `<GitHub Org>/<Repo Name>`. For `Bitbucket` repos, this will take the form of `<Workspace ID>/<Repo Name>`.
- `repo_type` - (Optional) The repository type. Valid values: `GitHub` or `TfsGit` or `Bitbucket` or `GitHub Enterprise`. Defaults to `Github`. If `repo_type` is `GitHubEnterprise`, must use existing project and GitHub Enterprise service connection.
- `service_connection_id` - (Optional) The service connection ID. Used if the `repo_type` is `GitHub` or `GitHubEnterprise`.
- `yml_path` - (Required) The path of the Yaml file describing the build definition.
- `github_enterprise_url` - (Optional) The Github Enterprise URL. Used if `repo_type` is `GithubEnterprise`.

`ci_trigger` block supports the following:

- `use_yaml` - (Optional) Use the azure-pipeline file for the build configuration. Defaults to `false`.
- `override` - (Optional) Override the azure-pipeline file and use a this configuration for all builds.

`ci_trigger` `override` block supports the following:

- `batch` - (Optional) If you set batch to true, when a pipeline is running, the system waits until the run is completed, then starts another run with all changes that have not yet been built. Defaults to `true`.
- `branch_filter` - (Optional) The branches to include and exclude from the trigger.
- `path_filter` - (Optional) Specify file paths to include or exclude. Note that the wildcard syntax is different between branches/tags and file paths.
- `max_concurrent_builds_per_branch` - (Optional) The number of max builds per branch. Defaults to `1`.
- `polling_interval` - (Optional) How often the external repository is polled. Defaults to `0`.
- `polling_job_id` - (Computed) This is the ID of the polling job that polls the external repository. Once the build definition is saved/updated, this value is set.

`pull_request_trigger` block supports the following:

- `use_yaml` - (Optional) Use the azure-pipeline file for the build configuration. Defaults to `false`.
- `initial_branch` - (Optional) When use_yaml is true set this to the name of the branch that the azure-pipelines.yml exists on. Defaults to `Managed by Terraform`.
- `forks` - (Required) Set permissions for Forked repositories.
- `override` - (Optional) Override the azure-pipeline file and use a this configuration for all builds.

`forks` block supports the following:

- `enabled` - (Required) Build pull requests form forms of this repository.
- `share_secrets` - (Required) Make secrets available to builds of forks.

`pull_request_trigger` `override` block supports the following:

- `auto_cancel` - (Optional) . Defaults to `true`.
- `branch_filter` - (Optional) The branches to include and exclude from the trigger.
- `path_filter` - (Optional) Specify file paths to include or exclude. Note that the wildcard syntax is different between branches/tags and file paths.

- `branch_filter` block supports the following:

- `include` - (Optional) List of branch patterns to include.
- `exclude` - (Optional) List of branch patterns to exclude.

- `path_filter` block supports the following:

- `include` - (Optional) List of path patterns to include.
- `exclude` - (Optional) List of path patterns to exclude.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the build definition
- `revision` - The revision of the build definition

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Build Definitions](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/definitions?view=azure-devops-rest-5.1)

## Import

Azure DevOps Build Definitions can be imported using the project name/definitions Id or by the project Guid/definitions Id, e.g.

```sh
terraform import azuredevops_build_definition.build "Test Project"/10
```

or

```sh
terraform import azuredevops_build_definition.build 782a8123-1019-xxxx-xxxx-xxxxxxxx/10
```
