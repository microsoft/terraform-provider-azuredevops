---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_devops"
description: |-
  Manages a Azure DevOps Service Connection service endpoint within Azure DevOps project.
---

# azuredevops_serviceendpoint_devops

Manages a Azure DevOps Service Connection service endpoint within Azure DevOps. Allows triggering of delegate pipelines, monitoring execution and collecting and consolidating artifacts produced in the delegate pipelines (yaml block `task: RunPipelines@1`). More details on Marketplace page: [RunPipelines](https://marketplace.visualstudio.com/items?itemName=CSE-DevOps.RunPipelines)

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_devops" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "my-pipelines-service-connection"
  organization          = "my-organization-name"
  auth_personal {
    # Also can be set with AZDO_PERSONAL_ACCESS_TOKEN environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  description = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `organization` - (Required) The organization name used for `Organization Url` and `Release API Url` fields.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.
- `auth_personal` - (Optional) An `auth_personal` block as documented below. Allows connecting using a personal access token.

`auth_personal` block supports the following:

- `personal_access_token` - (Required) The Personal Access Token for Azure DevOps Pipeline. It also can be set with AZDO_PERSONAL_ACCESS_TOKEN environment variable.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)

## Import

Azure DevOps Service Endpoint can be imported using the `project id`, `service connection id` , e.g.

```sh
 terraform import azuredevops_serviceendpoint_devops.serviceendpoint projectID/00000000-0000-0000-0000-000000000000
```
