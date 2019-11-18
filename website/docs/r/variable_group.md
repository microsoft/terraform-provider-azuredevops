# azuredevops_variable_group
Manages variable groups within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name = "Test Project"
}

resource "azuredevops_variable_group" "variablegroup" {
  project_id  = azuredevops_project.project.id
  name        = "Test Variable Group"
  description = "Test Variable Group Description"

  variable {
    name  = "key"
    value = "value"
  }

  variable {
    name      = "Account Password"
    value     = "p@ssword123"
    is_secret = true
  }
}
```

## Arugument Reference

The following arguments are supported:

* `project_id` - (Required) The Project in which this Variable Group exists or will be created.
* `name` - (Required) The name of the Variable Group.
* `description` - (Optional) The description of the Variable Group.
* `variable` - (Optional) One or more `variable` blocks as documented below.

A `variable` block supports the following:

* `name` - (Required) The key value used for the variable. Must be unique within the Variable Group.
* `value` - (Optional) The value of the variable. If omitted, it will default to empty string.
* `is_secret` - (Optional) A boolean flag describing if the variable value is sensitive. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Variable Group returned after creation in Azure DevOps.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-5.1)

## Import
Azure DevOps Variable groups can be imported using the project name/variable group Id or by the project Guid id/variable group Id, e.g.
 
 ```
 terraform import azuredevops_project.project "Test Project"/10
 or
 terraform import azuredevops_project.project 782a8123-1019-xxxx-xxxx-xxxxxxxx/10
 ```

*Note that for secret variables, the import command retrieve blank value in the tfstate.*

## PAT Permissions Required

- **Variable Groups**: Read, Create, & Manage