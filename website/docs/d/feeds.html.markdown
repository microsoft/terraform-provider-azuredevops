---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_feeds"
description: |-
  Gets information about an existing Feeds.
---

# Data Source: azuredevops_feeds

Use this data source to access information about an existing Feeds.

## Example Usage

```hcl
data "azuredevops_feeds" "example" {

}

output "id" {
  value = data.azuredevops_feeds.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Optional) The ID of the project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Feeds.

* `feeds` - A `feeds` block as defined below.

---

A `feeds` block exports the following:

* `badges_enabled` - If set, this feed supports generation of package badges.
* `description` - A description for the feed. Descriptions must not exceed 255 characters.
* `hide_deleted_package_versions` - If set, the feed will hide all deleted/unpublished versions
* `upstream_enabled` - This should always be true. Setting to false will override all sources in UpstreamSources.
* `upstream_sources` - A list of sources that this feed will fetch packages from. An empty list indicates that this feed will not search any additional sources for packages.


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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Feeds.
