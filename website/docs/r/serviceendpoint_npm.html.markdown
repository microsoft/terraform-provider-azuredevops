---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_npm"
description: |-
  Manages a npm server endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_npm

Manages a npm service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_npm" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example npm"
  url                   = "https://registry.npmjs.org"
  access_token          = "00000000-0000-0000-0000-000000000000"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `url` - (Required) URL of the npm registry to connect with.
- `access_token` - (Required) The access token for npm registry.
- `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)
- [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
- [npm User Token](https://docs.npmjs.com/about-access-tokens)

## Import

Azure DevOps Service Endpoint npm can be imported using the **projectID/serviceEndpointID**, e.g.

```shell
terraform import azuredevops_serviceendpoint_npm.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
