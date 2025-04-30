---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_black_duck"
description: |-
  Manages a Black Duck Detect service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_black_duck

Manages a Black Duck service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Black Duck Detect](https://marketplace.visualstudio.com/items?itemName=blackduck.blackduck-detect)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_black_duck" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://blackduck.com/"
  api_token             = "ffffffffffffffffff"
  service_endpoint_name = "Example Black Duck"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `server_url` - (Required) The server URL of the Black Duck Detect.

* `api_token` - (Required) The API token of the Black Duck Detect.

---

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Black Duck Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Black Duck Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Black Duck Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Black Duck Service Endpoint.

## Import

Azure DevOps Black Duck Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_black_duck.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
