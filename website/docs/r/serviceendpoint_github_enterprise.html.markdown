---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_github_enterprise"
description: |-
  Manages a GitHub-Enterprise-Server service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_github_enterprise

Manages a GitHub Enterprise Server service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_github_enterprise" "serviceendpoint_ghes_1" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample GitHub Enterprise"
  url                   = "https://github.contoso.com"
  description           = "Managed by Terraform"
  
  auth_personal {
    # Also can be set with AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `url` - (Required) Github Enterprise Server Url.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.
- `auth_personal` - (Optional) An `auth_personal` block as documented below. Allows connecting using a personal access token.

**NOTE: Github Apps can not be created or updated via terraform. You must install and configure the app on Github and then import it. You must also set the `description` to "" explicitly."**

`auth_personal` block supports the following:

- `personal_access_token` - (Required) The Personal Access Token for Github.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)

## Import

Azure DevOps Service Endpoint GitHub Enterprise Server can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
$ terraform import azuredevops_serviceendpoint_github_enterprise.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
