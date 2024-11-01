---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_generic"
description: |-
  Manages a generic service endpoint within Azure DevOps, which can be used to authenticate to any external server using
  basic authentication via a username and password.
---

# azuredevops_serviceendpoint_generic

Manages a generic service endpoint within Azure DevOps, which can be used to authenticate to any external server using
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

resource "azuredevops_serviceendpoint_generic" "example" {
  project_id            = azuredevops_project.example.id
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
- `server_url` - (Required) The URL of the server associated with the service endpoint.
- `username` - (Optional) The username used to authenticate to the server url using basic authentication.
- `password` - (Optional) The password or token key used to authenticate to the server url using basic authentication.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The name of the service endpoint.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Import

Azure DevOps Service Endpoint Generic can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_generic.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
