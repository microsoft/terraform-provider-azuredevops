# azuredevops_project
Manages a project within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}
```

## Arugument Reference

The following arguments are supported:

* `project_name` - (Required) The Project Name.
* `description` - (Optional) The Description of the Project.
* `visibility` - (Optional) Specifies the visibility of the Project. Possible values are `private` or `public`. - private is the default.
* `version_control` - (Optional) Specifies the version control system. Possible values are `Git` or `Tfvc`. - Git is the default. If you change this value on update, terraform will re-create the project.
* `work_item_template` - (Optional) Specifies the work item template. - Agile is the default. If you change this value on update, terraform will re-create the project.

## Attributes Reference

The following attributes are exported:

* `id` - The Project ID of the Project.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Projects](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects?view=azure-devops-rest-5.1)

## Import
Azure DevOps Projects can be imported using the project name or by the project Guid id, e.g.

```
terraform import azuredevops_project.project "Test Project"
or
terraform import azuredevops_project.project 782a8123-1019-xxxx-xxxx-xxxxxxxx
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage