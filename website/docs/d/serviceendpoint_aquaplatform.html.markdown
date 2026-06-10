---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_aquaplatform"
description: |-
  Gets information about an existing Aqua Platform Service Endpoint.
---

# Data Source : azuredevops_serviceendpoint_aquaplatform

Use this data source to access information about an existing Aqua Platform Service Endpoint.

## Example Usage

```hcl
data "azuredevops_serviceendpoint_aquaplatform" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Aqua Platform"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_aquaplatform.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

---

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `id` - The ID of the Aqua Platform Service Endpoint.

* `authorization` - The Authorization scheme.

* `aqua_platform_url` - The URL of the Aqua Platform.

* `aqua_auth_url` - The URL used for authentication.

* `description` - The description of the Service Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Aqua Platform Service Endpoint.
