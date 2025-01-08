---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_checkmarx_sca"
description: |-
  Manages a Checkmarx SCA service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_checkmarx_sca

Manages a Checkmarx SCA service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Checkmarx SAST](https://marketplace.visualstudio.com/items?itemName=checkmarx.cxsast)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_checkmarx_sca" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Checkmarx SCA"
  access_control_url    = "https://accesscontrol.com"
  server_url            = "https://server.com"
  web_app_url           = "https://webapp.com"
  account               = "account"
  username              = "username"
  password              = "password"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `access_control_url` - (Required) The Access Control URL of the Checkmarx SCA.

* `server_url` - (Required) The Server URL of the Checkmarx SCA.

* `web_app_url` - (Required) The Web App URL of the Checkmarx SCA.
  
* `account` - (Required) The account of the Checkmarx SCA.

* `username` - (Required) The username of the Checkmarx SCA.

* `password` - (Required) The password of the Checkmarx SCA.

* `team` - (Optional) The full team name of the Checkmarx.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Checkmarx SCA Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Checkmarx SCA Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Checkmarx SCA Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Checkmarx SCA Service Endpoint.

## Import

Azure DevOps Service Endpoint Check Marx SCA can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_checkmarx_sca.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
