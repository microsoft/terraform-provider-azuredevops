# azuredevops_azure_git_repository
Manages a git repository within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_azure_git_repository" "repo" {
  project_id = azuredevops_project.project.id
  name       = "Sample Empty Git Repository"
  initialization {
    init_type   = "Clean"
  }
```


```hcl
resource "azuredevops_azure_git_repository" "repo" {
  project_id = azuredevops_project.project.id
  name       = "Sample Fork an Existing Repository"
  initialization {
    init_type   = "Fork"
    source_type = ""
    source_url  = ""
  }
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `name` - (Required) The name of the git repository.
* `initialization` - (Required) An `initialization` block as documented below.

`initialization` block supports the following:

* `init_type` - (Required) The type of repository to create. Valid values: `Uninitialized`, `Clean`, `Fork`, or `Import`. Defaults to `Uninitialized`.
* `source_type` - (Optional) Type type of the source repository. Used if the init type is `Fork` or `Import`.
* `source_url` - (Optional) The url of the source repository. Used if the init type is `Fork` or `Import`.

## Attributes Reference

In addition to all arguments above, except `initialization`, the following attributes are exported:

* `id` - The ID of the agent pool.

* `default_branch` - The name of the default branch.
* `is_fork` - True if the repository was created as a fork.
* `remote_url` - If `init_type` is `Fork` the url of the remote repository.
* `size` - Size in bytes.
* `ssh_url` - Git SSH Url of the repository.
* `url` - Git Url of the repository.
* `web_url` - Web link to the repository.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/repositories?view=azure-devops-rest-5.1)
