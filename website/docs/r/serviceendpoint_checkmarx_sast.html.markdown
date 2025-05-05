---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_checkmarx_sast"
description: |-
  Manages a Checkmarx SAST service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_checkmarx_sast

Manages a Checkmarx SAST service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Checkmarx SAST](https://marketplace.visualstudio.com/items?itemName=checkmarx.cxsast)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_checkmarx_sast" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Checkmarx SAST"
  server_url            = "https://server.com"
  username              = "username"
  password              = "password"
  team                  = "team"
  preset                = "preset"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `server_url` - (Required) The Server URL of the Checkmarx SAST.

* `username` - (Required) The username of the Checkmarx SAST.

* `password` - (Required) The password of the Checkmarx SAST.

---

* `team` - (Optional) The full team name of the Checkmarx.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

* `preset` - (Optional) Predefined sets of queries that you can select when Creating, Configuring and Branching Projects. Predefined presets are provided by Checkmarx and you can configure your own. You can also import and export presets (on the server).In Service Connection if preset(optional) value is added, then it will igonres Preset available in pipeline and uses preset available in service connection only.If Preset is blank in service connection then it will use pipelines preset.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Checkmarx SAST Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Checkmarx SAST Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Checkmarx SAST Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Checkmarx SAST Service Endpoint.

## Import

Azure DevOps Service Endpoint Check Marx SAST can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_checkmarx_sast.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
