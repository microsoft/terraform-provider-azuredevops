---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_serviceendpoint_generic"
description: |-
  Gets information about an existing Generic Service Endpoint.
---

# Data Source: azuredevops_serviceendpoint_generic

Use this data source to access information about an existing Generic Service Endpoint.

## Example Usage

### By Service Endpoint ID

```hcl
data "azuredevops_project" "sample" {
  name = "Sample Project"
}

data "azuredevops_serviceendpoint_generic" "serviceendpoint" {
  project_id          = data.azuredevops_project.sample.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

output "service_endpoint_name" {
  value = data.azuredevops_serviceendpoint_generic.serviceendpoint.service_endpoint_name
}
```

### By Service Endpoint Name

```hcl
data "azuredevops_project" "sample" {
  name = "Sample Project"
}

data "azuredevops_serviceendpoint_generic" "serviceendpoint" {
  project_id            = data.azuredevops_project.sample.id
  service_endpoint_name = "Example-Service-Endpoint"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_generic.serviceendpoint.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

---

* `service_endpoint_id` - (Optional) The ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) Service Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Generic Service Endpoint.

* `authorization` - A `authorization` block as defined below.

* `description` - The description of the Service Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Generic Service Endpoint.
