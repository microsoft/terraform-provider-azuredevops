---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_custom"
description: |-
  Manages a custom service endpoint within Azure DevOps, which can be used to authenticate to any external server using
  basic authentication via a username and password.
---

# azuredevops_serviceendpoint_custom

Manages a custom service endpoint within Azure DevOps, which can be used to authenticate to any external server using
basic authentication via a username and password.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_custom" "example" {
  project_id            = azuredevops_project.example.id
  service_type          = "custom type name"
  server_url            = "https://some-server.example.com"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Example Generic"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The service endpoint name.
- `service_type` - (Required) The Service Type of the server associated with the service endpoint.
- `server_url` - (Required) The URL of the server associated with the service endpoint.
- `username` - (Optional) The username used to authenticate to the server url using basic authentication.
- `password` - (Optional) The password or token key used to authenticate to the server url using basic authentication.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

Obs.: Access [Azure DevOps Service REST API 7.1 - Service Endpoints Types List](https://dev.azure.com/{organization}/_apis/serviceendpoint/types?api-version=7.1-preview.1) to get the Type.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The name of the service endpoint.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)
- [Azure DevOps Service REST API 7.1 - Service Endpoints Types List](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/types/list?view=azure-devops-rest-7.1&tabs=HTTP)

## Import

Azure DevOps Service Endpoint Generic can be imported using **projectID/serviceEndpointID** or
**projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_custom.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
