---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_serviceendpoint_dockerregistry"
description: |-
  Gets information about an existing Docker Registry Service Endpoint.
---

# Data Source : azuredevops_serviceendpoint_dockerregistry

Use this data source to access information about an existing Docker Registry Service Endpoint.

## Example Usage

### By Service Endpoint ID

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_serviceendpoint_dockerregistry" "example" {
  project_id          = data.azuredevops_project.example.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

output "service_endpoint_name" {
  value = data.azuredevops_serviceendpoint_dockerregistry.example.service_endpoint_name
}
```

### By Service Endpoint Name

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_serviceendpoint_dockerregistry" "example" {
  project_id            = data.azuredevops_project.example.id
  service_endpoint_name = "Example-Service-Endpoint"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_dockerregistry.serviceendpoint.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - the ID of the Service Endpoint.

* `service_endpoint_name` - the Name of the Service Endpoint.

~> **NOTE:** 1. One of either `service_endpoint_id` or `service_endpoint_name` must be specified.
    <br>2. When supplying `service_endpoint_name`, take care to ensure that this is a unique name.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `id` - The ID of the Azure Resource Manager Service Endpoint.

* `authorization` - The Authorization scheme.

* `description` - The Service Endpoint description.

* `docker_registry` - The URL of the Docker registry.

* `docker_username` - The identifier of the Docker account user.

* `docker_email` - The email for Docker account user.

* `docker_password` - The password for the account user identified above.

* `registry_type` - Can be "DockerHub" or "Others" (Default "DockerHub")

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Azure Resource Manager Service Endpoint.
