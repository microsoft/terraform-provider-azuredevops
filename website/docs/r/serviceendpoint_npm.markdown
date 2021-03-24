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
resource "azuredevops_project" "project" {
  name               = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_npm" "serviceendpoint" {

  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample npm"
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
- `access_token` - (Required) Authentication Token generated through npm (go to My Account > Security > Generate Tokens).
- `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
- [npm User Token](https://docs.npmjs.com/integrations/integrating-npm-with-external-services)

## Import

Azure DevOps Service Endpoint npm can be imported using the **projectID/serviceEndpointID**, e.g.

```shell
$ terraform import azuredevops_serviceendpoint_npm.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
