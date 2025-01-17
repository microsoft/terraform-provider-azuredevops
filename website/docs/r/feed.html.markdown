---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_feed"
description: |-
  Manages Feed within Azure DevOps organization.
---

# azuredevops_feed

Manages Feed within Azure DevOps organization.

## Example Usage

### Create Feed in the scope of whole Organization
```hcl
resource "azuredevops_feed" "example" {
  name = "examplefeed"
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
  name       = "examplefeed"
  project_id = azuredevops_project.example.id
}
```

### Create Feed with Soft Delete
```hcl
resource "azuredevops_feed" "example" {
  name = "examplefeed"
  features {
    permanent_delete = false
  }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Feed.

---

* `project_id` - (Optional) The ID of the Project Feed is created in. If not specified, feed will be created at the organization level.

* `features`- (Optional) A `features` blocks as documented below.

~> **Note** *Because of ADO limitations feed name can be **reserved** for up to 15 minutes after permanent delete of the feed*

---

`features` block supports the following:

* `permanent_delete` - (Optional) Determines if Feed should be Permanently removed, Defaults to `false`
* `restore` - (Optional) Determines if Feed should be Restored during creation (if possible), Defaults to `false`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Feed.
* `name` - The name of the Feed.
* `project_id` - The ID of the Project Feed is created in (if one exists).

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Feed.
* `read` - (Defaults to 5 minute) Used when retrieving the Feed.
* `update` - (Defaults to 10 minutes) Used when updating the Feed.
* `delete` - (Defaults to 10 minutes) Used when deleting the Feed.

## Import

Azure DevOps Feed can be imported using the Project ID and Feed ID or Feed ID e.g.:

```sh
terraform import azuredevops_feed.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

or 

```sh
terraform import azuredevops_feed.example 00000000-0000-0000-0000-000000000000
```

