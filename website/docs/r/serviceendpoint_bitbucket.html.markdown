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
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_bitbucket" "example" {
  project_id            = azuredevops_project.example.id
  email                 = "email@example.com"
  api_token             = "api_token"
  service_endpoint_name = "Example Bitbucket"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

---

* `email` - (Optional) Bitbucket account email. Used together with `api_token` to authenticate using an Atlassian API token.

* `api_token` - (Optional) Bitbucket account API token. Used together with `email` to authenticate using an Atlassian API token.

* `username` - (Optional) Bitbucket account username. Used together with `password` to authenticate using an app password. **Deprecated**: Bitbucket Cloud has deprecated app password (username and password) authentication. Use `email` and `api_token` instead.

* `password` - (Optional) Bitbucket account password. Used together with `username` to authenticate using an app password. **Deprecated**: Bitbucket Cloud has deprecated app password (username and password) authentication. Use `email` and `api_token` instead.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

~> **NOTE:** Exactly one authentication method must be configured. Provide either `email` + `api_token` (recommended) or `username` + `password` (deprecated).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Bitbucket Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Bitbucket Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Bitbucket Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Bitbucket Service Endpoint.

## Import

Azure DevOps Bitbucket Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_bitbucket.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
