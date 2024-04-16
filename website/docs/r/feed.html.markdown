---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_feed"
description: |-
  Manages creation of the Feed within Azure DevOps organization.
---

# Data Source: azuredevops_feed

Manages creation of the Feed within Azure DevOps organization.

## Example Usage

### Create Feed in the scope of whole Organization
```hcl
resource "azuredevops_feed" "example" {
  name = "releases"
}
```

### Create Feed in the scope of a Project
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_feed" "example" {
  name       = "releases"
  project_id = azuredevops_project.example.id
}
```

### Create Feed with Soft Delete
```hcl
resource "azuredevops_feed" "example" {
  name             = "releases"
  permanent_delete = false
}
```


## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Feed. *Because of ADO limitations feed name can be **reserved** for up to 15 minutes after permanent delete of the feed*
- `project_id` - (Optional) The ID of the Project Feed is created in. If not specified, feed will be created at the organization level.
- `permanent_delete` - (Optional) Determines if Feed should be Permanently removed, default value is `true`

## Attributes Reference

The following attributes are exported:

- `name` - The name of the Feed.
- `project_id` - The ID of the Project Feed is created in (if one exists).
- `restored` - Determines if Feed was restored after Soft Delete

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)