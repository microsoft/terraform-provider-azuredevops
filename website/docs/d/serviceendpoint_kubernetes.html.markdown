---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_kubernetes"
description: |-
  Gets information about an existing Kubernetes Service Endpoint. 
---

# Data Source : azuredevops_serviceendpoint_kubernetes

Use this data source to access information about an existing Kubernetes Service Endpoint.

## Example Usage

```hcl
data "azuredevops_serviceendpoint_kubernetes" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Kubernetes"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_kubernetes.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.


## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `authorization` - Specifies the Authorization Scheme Map.
* `description` - Specifies the description of the Service Endpoint.
