---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_type"
description: |-
  Gets information about a specific service endpoint type in Azure DevOps.
---

# Data Source: azuredevops_serviceendpoint_type

Use this data source to query a specific service endpoint type and retrieve its parameters and authorization schemes.

## Example Usage

### Basic usage - Get service endpoint type information

```hcl
data "azuredevops_serviceendpoint_type" "generic" {
  name = "generic"
}

output "generic_type_info" {
  value = {
    display_name = data.azuredevops_serviceendpoint_type.generic.display_name
    auth_schemes = data.azuredevops_serviceendpoint_type.generic.authentication_schemes
    parameters   = data.azuredevops_serviceendpoint_type.generic.parameters
  }
}
```

### Get authorization parameters for a specific auth scheme

```hcl
data "azuredevops_serviceendpoint_type" "generic" {
  name                 = "generic"
  authorization_scheme = "UsernamePassword"
}

output "auth_parameters" {
  value = data.azuredevops_serviceendpoint_type.generic.authorization_parameters
}
```

### Use in service endpoint resource

```hcl
data "azuredevops_serviceendpoint_type" "bitbucket" {
  name                 = "bitbucket"
  authorization_scheme = "UsernamePassword"
}

resource "azuredevops_project" "example" {
  name = "Example Project"
}

# Use the data source to understand available parameters
resource "azuredevops_serviceendpoint_generic_v2" "example" {
  project_id           = azuredevops_project.example.id
  name                 = "Example Bitbucket"
  description          = "Managed by Terraform"
  type                 = data.azuredevops_serviceendpoint_type.bitbucket.name
  server_url           = "https://bitbucket.org"
  authorization_scheme = data.azuredevops_serviceendpoint_type.bitbucket.authorization_scheme

  authorization_parameters = {
    username = "my-username"
    password = "my-password"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service endpoint type to query.
* `authorization_scheme` - (Optional) The authorization scheme to retrieve parameters for. When provided, the `authorization_parameters` output will be populated with the parameters specific to this authentication scheme.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint type (typically same as name).
* `display_name` - The display name of the service endpoint type.
* `description` - The description of the service endpoint type.
* `ui_contribution_id` - The UI contribution ID for this service endpoint type.
* `authentication_schemes` - A list of available authentication schemes for this service endpoint type.
* `parameters` - A map of parameter names to their default values. These are the parameters that can be set in the `parameters` block of a service endpoint resource.
* `authorization_parameters` - A map of authorization parameter names to their default values. This is only populated when `authorization_scheme` is provided. These are the parameters that can be set in the `authorization_parameters` block of a service endpoint resource.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoint Types](https://learn.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/types/list)

