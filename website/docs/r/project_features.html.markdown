---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_features"
description: |-
  Manages features for Azure DevOps projects.
---

# azuredevops_project_features

Manages features for Azure DevOps projects

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project_features" "example-features" {
  project_id = azuredevops_project.example.id
  features = {
    testplans = "disabled"
    artifacts = "enabled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `projectd_id` - (Required) The `id` of the project for which the project features will be managed.

* `features` - (Required) Defines the status (`enabled`, `disabled`) of the project features.  Valid features `boards`, `repositories`, `pipelines`, `testplans`, `artifacts`

  | Features     | Possible Values   |
  |--------------|-------------------|
  | boards       | enabled, disabled |
  | repositories | enabled, disabled |
  | pipelines    | enabled, disabled |
  | testplans    | enabled, disabled |
  | artifacts    | enabled, disabled |

  ~> **NOTE:** It's possible to define project features both within the [`azuredevops_project_features` resource](project_features.html) and 
    via the `features` block by using the [`azuredevops_project` resource](project.html). 
    However it's not possible to use both methods to manage features, since there'll be conflicts.

## Attributes Reference

No attributes are exported

## Relevant Links

No official documentation available

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Project Features.
* `read` - (Defaults to 5 minute) Used when retrieving the Project Features.
* `update` - (Defaults to 10 minutes) Used when updating the Project Features.
* `delete` - (Defaults to 10 minutes) Used when deleting the Project Features.
 
## Import

Azure DevOps feature settings can be imported using the project id, e.g.

```sh
terraform import azuredevops_project_features.example 00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage
