---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azuredevops"
description: |-
  Manages a Azure DevOps service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azuredevops

Manages an Azure DevOps service endpoint within Azure DevOps.

~> **Note** Prerequisite: Extension [Configurable Pipeline Runner](https://marketplace.visualstudio.com/items?itemName=CSE-DevOps.RunPipelines) has been installed for the organization. 

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azuredevops" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Azure DevOps"
  org_url               = "https://dev.azure.com/testorganization"
  release_api_url       = "https://vsrm.dev.azure.com/testorganization"
  personal_access_token = "0000000000000000000000000000000000000000000000000000"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `org_url` - (Required) The organization URL.
- `release_api_url` - (Required) The URL of the release API.
- `personal_access_token` - (Required) The Azure DevOps personal access token.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint Azure DevOps can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_azuredevops.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
