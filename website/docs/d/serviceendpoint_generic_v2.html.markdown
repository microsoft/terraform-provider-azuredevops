---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_generic_v2"
description: |-
  Gets information about a Generic Service Endpoint (v2) within Azure DevOps.
---

# Data Source: azuredevops_serviceendpoint_generic_v2

Use this data source to access information about an existing Generic Service Endpoint (v2) within Azure DevOps.

## Example Usage

```hcl
# Get service endpoint by ID
data "azuredevops_serviceendpoint_generic_v2" "example" {
  project_id          = azuredevops_project.example.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

# Get service endpoint by name
data "azuredevops_serviceendpoint_generic_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Generic Service Endpoint"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to which the service endpoint belongs.
* `id` - (Optional) The ID of the service endpoint to retrieve. One of `service_endpoint_id` or `service_endpoint_name` must be specified.
* `name` - (Optional) The name of the service endpoint to retrieve. One of `service_endpoint_id` or `service_endpoint_name` must be specified.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of the service endpoint.
* `type` - The type of the service endpoint.
* `description` - The description of the service endpoint.
* `server_url` - The URL of the server associated with the service endpoint.
* `authorization_scheme` - The authorization type of the service endpoint.
* `authorization_parameters` - The authorization parameters of the service endpoint.
* `data` - Additional data of the service endpoint.
