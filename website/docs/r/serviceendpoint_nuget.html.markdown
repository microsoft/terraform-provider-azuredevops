---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_nuget"
description: |-
  Manages a NuGet service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_nuget

Manages a NuGet service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_nuget" "example" {
  project_id            = azuredevops_project.example.id
  api_key               = "apikey"
  service_endpoint_name = "Example NuGet"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `feed_url` - (Required) The URL for the feed. This will generally end with `index.json`.

---
- `api_key` - (Optional) The API Key used to connect to the endpoint.
- `personal_access_token` - (Optional) The Personal access token used to  connect to the endpoint. Personal access tokens are applicable only for NuGet feeds hosted on other Azure DevOps Services organizations or Azure DevOps Server 2019 (or later).
- `username` - (Optional) The account username used to connect to the endpoint.
- `password` - (Optional) The account password used to connect to the endpoint

~> **Note** Only one of `api_key` or `personal_access_token` or  `username`, `password` can be set at the same time.

- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the NuGet Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the NuGet Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the NuGet Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the NuGet Service Endpoint.

## Import

Azure DevOps NuGet Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_nuget.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
