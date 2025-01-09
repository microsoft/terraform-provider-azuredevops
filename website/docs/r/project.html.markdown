---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project"
description: |-
  Manages a project within Azure DevOps organization.
---

# azuredevops_project

Manages a project within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
  features = {
     testplans = "disabled"
     artifacts = "disabled"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Project Name.
- `description` - (Optional) The Description of the Project.
- `visibility` - (Optional) Specifies the visibility of the Project. Valid values: `private` or `public`. Defaults to `private`.
- `version_control` - (Optional) Specifies the version control system. Valid values: `Git` or `Tfvc`. Defaults to `Git`.
- `work_item_template` - (Optional) Specifies the work item template. Valid values: `Agile`, `Basic`, `CMMI`, `Scrum` or a custom, pre-existing one. Defaults to `Agile`. An empty string will use the parent organization default.
- `features` - (Optional) Defines the status (`enabled`, `disabled`) of the project features.
   Valid features are `boards`, `repositories`, `pipelines`, `testplans`, `artifacts`

> **NOTE:**
> It's possible to define project features both within the [`azuredevops_project_features` resource](project_features.html) and
> via the `features` block by using the [`azuredevops_project` resource](project.html).
> However it's not possible to use both methods to manage features, since there'll be conflicts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The Project ID of the Project.
- `process_template_id` - The Process Template ID used by the Project.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Projects](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Project.
* `read` - (Defaults to 5 minute) Used when retrieving the Project.
* `update` - (Defaults to 10 minutes) Used when updating the Project.
* `delete` - (Defaults to 10 minutes) Used when deleting the Project.

## Import

Azure DevOps Projects can be imported using the project name or by the project Guid, e.g.

```sh
terraform import azuredevops_project.example "Example Project"
```

or

```sh
terraform import azuredevops_project.example 00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage
- **Work Items**: Read