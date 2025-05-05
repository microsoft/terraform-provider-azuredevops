---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_checkmarx_one"
description: |-
  Manages a Checkmarx One service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_checkmarx_one

Manages a Checkmarx One service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Checkmarx AST](https://marketplace.visualstudio.com/items?itemName=checkmarx.checkmarx-ast-azure-plugin)

## Example Usage

### Authorize with API Key

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Checkmarx One"
  server_url            = "https://server.com"
  api_key               = "apikey"
}
```

### Authorize with Client ID and Secret

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_checkmarx_one" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Checkmarx One"
  server_url            = "https://server.com"
  client_id             = "clientid"
  client_secret         = "secret"
  authorization_url     = "https://authurl.com"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `server_url` - (Required) The Server URL of the Checkmarx One Service.

---

* `authorization_url` - (Optional) The URL of Checkmarx Authorization. Used when using `client_id` and `client_secret` authorization.

* `api_key` - (Optional) The account of the Checkmarx One. Conflict with `client_id` and `client_secret`.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

* `client_id` - (Optional) The Client ID of the Checkmarx One. Conflict with `api_key`

* `client_secret` - (Optional) The Client Secret of the Checkmarx One. Conflict with `api_key`

~> **Note** At least one of `api_key` and `client_id`, `client_secret` must be set

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Checkmarx One Service.
* `read` - (Defaults to 1 minute) Used when retrieving the Checkmarx One Service.
* `update` - (Defaults to 2 minutes) Used when updating the Checkmarx One Service.
* `delete` - (Defaults to 2 minutes) Used when deleting the Checkmarx One Service.

## Import

Azure DevOps Service Endpoint Check Marx One can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_checkmarx_one.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
