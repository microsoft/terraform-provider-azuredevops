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
provider "azuredevops" {
  version = ">= 0.0.1"
}

data "azuredevops_project" "tf-project-test-001" {
  project_name = "Test Project"
}

resource "azuredevops_project_features" "my-project-features" {
  project_id = data.azuredevops_project.tf-project-test-001.id
  features = {
      "testplans" = "disabled"
      "artifacts" = "enabled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `projectd_id` - (Required) The `id` of the project for which the project features will be managed.
* `features` - (Required) Defines the status (`enabled`, `disabled`) of the project features.  
   Valid features `boards`, `repositories`, `pipelines`, `testplans`, `artifacts`

> **NOTE:**  
> It's possible to define project features both within the [`azuredevops_project_features` resource](project_features.html) and 
> via the `features` block by using the [`azuredevops_project` resource](project.html).
> However it's not possible to use both methods to manage group members, since there'll be conflicts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

## Relevant Links

No official documentation available

## Import

Azure DevOps feature settings can be imported using the project id, e.g.

```sh
$ terraform import azuredevops_project_features.project_id 2785562e-8f45-4534-a10e-b9ca1666b17e
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage
