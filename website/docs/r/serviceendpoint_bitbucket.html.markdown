---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_bitbucket"
description: |-
  Manages a Bitbucket service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_bitbucket

Manages a Bitbucket service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_bitbucket" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  username              = "username"
  password              = "password"
  service_endpoint_name = "Sample Bitbucket"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `username` - (Required) Bitbucket account username.
- `password` - (Required) Bitbucket account password.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint Bitbucket can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_bitbucket.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
