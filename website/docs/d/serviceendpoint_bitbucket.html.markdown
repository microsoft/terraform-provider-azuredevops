---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_serviceendpoint_bitbucket"
description: |-
  Gets information about an existing Bitbucket Service Endpoint.
---

# Data Source : azuredevops_serviceendpoint_bitbucket

Use this data source to access information about an existing Bitbucket service Endpoint.

## Example Usage

### By Service Endpoint ID

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_serviceendpoint_bitbucket" "example" {
  project_id          = data.azuredevops_project.example.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

output "service_endpoint_name" {
  value = data.azuredevops_serviceendpoint_bitbucket.example.service_endpoint_name
}
```

### By Service Endpoint Name

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_serviceendpoint_bitbucket" "example" {
  project_id            = data.azuredevops_project.example.id
  service_endpoint_name = "Example"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_bitbucket.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.
~> **NOTE:** When supplying `service_endpoint_name`, take care to ensure that this is a unique name.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `authorization` - Specifies the Authorization Scheme Map.
* `description` - Specifies the description of the Service Endpoint.

## PAT Permissions Required

- **vso.serviceendpoint**: Grants the ability to read service endpoints.
