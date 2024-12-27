layout: "azuredevops"
---
page_title: "AzureDevops: azuredevops_feed_retention_policy"
description: |-
  Manages the Feed Retention Policy within Azure DevOps organization.
---

# azuredevops_feed_retention_policy

Manages the Feed Retention Policy within Azure DevOps.

## Example Usage - Project Feed
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_feed" "example" {
  name       = "ExampleFeed"
  project_id = azuredevops_project.example.id
}

resource "azuredevops_feed_retention_policy" "example" {
  project_id                                = azuredevops_project.example.id
  feed_id                                   = azuredevops_feed.example.id
  count_limit                               = 20
  days_to_keep_recently_downloaded_packages = 30
}
```

## Example Usage - Organization Feed
```hcl
resource "azuredevops_feed" "example" {
  name       = "examplefeed"
}

resource "azuredevops_feed_retention_policy" "example" {
  feed_id                                   = azuredevops_feed.example.id
  count_limit                               = 20
  days_to_keep_recently_downloaded_packages = 30
}
```

## Argument Reference

The following arguments are supported:

* `feed_id` - (Required) The ID of the Feed. Changing this forces a new resource to be created.

* `count_limit`- (Required) The maximum number of versions per package.

* `days_to_keep_recently_downloaded_packages`- (Required) The days to keep recently downloaded packages.

* `project_id` - (Optional) The ID of the Project. If not specified, Feed will be created at the organization level. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `feed_id` - The ID of the Feed.
* `project_id` - The ID of the Project.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)

## Import

Azure DevOps Feed Retention Policy can be imported using the Project ID and Feed ID or Feed ID e.g.:

```sh
terraform import azuredevops_feed_retention_policy.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

or 

```sh
terraform import azuredevops_feed_retention_policy.example 00000000-0000-0000-0000-000000000000
```

