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

### Create Feed in the scope of whole Organization
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_feed" "example" {
  name = "releases"
  project_id = azuredevops_project.example.id
}
```


## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Feed.
- `project_id` - (Optional) ID of the Project Feed is created in.


## Attributes Reference

The following attributes are exported:

- `name` - Name of the Feed.
- `project_id` - ID of the Project Feed is created in (if one exists).

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)