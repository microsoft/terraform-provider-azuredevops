---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_feed"
description: |-
  Use this data source to access information about existing Feed within a given project in Azure DevOps.
---

# Data Source: azuredevops_feed

Use this data source to access information about existing Feed within a given project in Azure DevOps.

## Example Usage

### Basic Example
```hcl
data "azuredevops_feed" "example" {
  name = "releases"
}
```

### Access feed within a project
```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_feed" "example" {
  name = "releases"
  project_id = data.azuredevops_project.example.id
}
```


## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Feed.
- `feed_id` - (Required) ID of the Feed.

~> **Note** Only one of `name` or `feed_id` can be set at the same time.

- `project_id` - (Optional) ID of the Project Feed is created in.


## Attributes Reference

The following attributes are exported:

- `name` - The name of the Feed.
- `feed_id` - The ID of the Feed.
- `project_id` - The ID of the Project.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management/get-feed?view=azure-devops-rest-7.0)