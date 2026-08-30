---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_feed"
description: |-
  Manages Feed within Azure DevOps organization.
---

# azuredevops_feed

Manages a Feed.

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

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Feed. Changing this forces a new Feed to be created.

---

* `features` - (Optional) One or more `features` blocks as defined below.

* `project_id` - (Optional) The ID of the project. Changing this forces a new Feed to be created.

* `upstream_sources` - (Optional) One or more `upstream_sources` blocks as defined below.

~> **Note** *Because of ADO limitations feed name can be **reserved** for up to 15 minutes after permanent delete of the feed*

---

A `features` block supports the following:

* `permanent_delete` - (Optional) Determines if Feed should be Permanently removed, Defaults to `false`
* `restore` - (Optional) Determines if Feed should be Restored during creation (if possible), Defaults to `false`

---

A `upstream_sources` block supports the following:

* `location` - (Required) Consistent locator for connecting to the upstream source.
* `name` - (Required) Display name.
* `protocol` - (Required) Package type associated with the upstream source.
* `internal_upstream_collection_id` - (Optional) For an internal upstream type, track the Azure DevOps organization that contains it.
* `internal_upstream_feed_id` - (Optional) For an internal upstream type, track the feed id being referenced.
* `internal_upstream_view_id` - (Optional) For an internal upstream type, track the project of the feed being referenced.
* `service_endpoint_id` - The identity of the service endpoint that holds credentials to use when accessing the upstream.
* `service_endpoint_project_id` - Specifies the projectId of the Service Endpoint.
* `upstream_source_type` - (Optional) Source type, such as public or internal.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Feed.
* `name` - The name of the Feed.
* `project_id` - The ID of the Project Feed is created in (if one exists).

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Feed.
* `read` - (Defaults to 5 minutes) Used when retrieving the Feed.
* `update` - (Defaults to 10 minutes) Used when updating the Feed.
* `delete` - (Defaults to 10 minutes) Used when deleting the Feed.

## Import

Feeds can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_feed.example 00000000-0000-0000-0000-000000000000
```
