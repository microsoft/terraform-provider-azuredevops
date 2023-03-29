---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_nuget"
description: |-
  Manages a NuGet server endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_nuget

Manages a NuGet service endpoint within Azure DevOps.

## Example Usage

One of the meta-argument blocks `authentication_token` or `authentication_basic` or `authentication_none` needs to be used.

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_nuget" "example_authentication_token" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example NuGet"
  description           = "Managed by Terraform"
  url                   = "https://api.nuget.org/v3/index.json"
  authentication_token {
    token = "AbcDEf123_0x"
  }
}

resource "azuredevops_serviceendpoint_nuget" "example_authentication_basic" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example NuGet"
  description           = "Managed by Terraform"
  url                   = "https://api.nuget.org/v3/index.json"
  authentication_basic {
    username = "username"
    password = "password"
  }
}

resource "azuredevops_serviceendpoint_nuget" "example_authentication_none" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example NuGet"
  description           = "Managed by Terraform"
  url                   = "https://api.nuget.org/v3/index.json"
  authentication_none {
    key = "AbcDEf123_0x"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `url` - (Required) URL for the feed. This will generally end with index.json. For nuget.org, use https://api.nuget.org/v3/index.json.
  _Note: URL should not end in a slash character._

* either `authentication_token` or `authentication_basic` or `authentication_none` (one is required)
  - `authentication_token`
    - `token` - Personal access tokens are applicable only for NuGet feeds hosted on other Azure DevOps Services organizations or Azure DevOps Server 2019 (or later).
  - `authentication_basic`
    - `username` - Username for connecting to the endpoint.
    - `password` - Password for connecting to the endpoint.
  - `authentication_none`
    - `key` - ApiKey (only for push).
* `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)
- [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
- [NuGet ApiKey](https://learn.microsoft.com/en-in/nuget/nuget-org/scoped-api-keys)
- [NuGet packages in Azure Artifacts](https://learn.microsoft.com/en-us/azure/devops/artifacts/get-started-nuget?view=azure-devops&tabs=windows)

## Import

Azure DevOps Service Endpoint NuGet can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_nuget.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
