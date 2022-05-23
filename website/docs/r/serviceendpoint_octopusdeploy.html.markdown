---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_octopusdeploy"
description: |-
  Manages an Octopus Deploy service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_octopusdeploy

Manages an Octopus Deploy service endpoint within Azure DevOps. Using this service endpoint requires you to install [Octopus Deploy](https://marketplace.visualstudio.com/items?itemName=octopusdeploy.octopus-deploy-build-release-tasks).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_octopusdeploy" "example" {
  project_id            = azuredevops_project.example.id
  url                   = "https://octopus.com"
  api_key               = "000000000000000000000000000000000000"
  service_endpoint_name = "Example Octopus Deploy"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `url` - (Required) Octopus Server url.
- `api_key` - (Required) API key to connect to Octopus Deploy.
- `ignore_ssl_error` - (Optional) Whether to ignore SSL errors when connecting to the Octopus server from the agent. Default to `false`.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint Octopus Deploy can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_octopusdeploy.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
