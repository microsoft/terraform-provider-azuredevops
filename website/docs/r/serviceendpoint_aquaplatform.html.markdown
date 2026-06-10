---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_aquaplatform"
description: |-
  Manages an Aqua Platform service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_aquaplatform
Manages an Aqua Platform service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_aquaplatform" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Aqua Platform"
  aqua_platform_url     = "https://aqua.com"
  aqua_key              = "00000000-0000-0000-0000-000000000000"
  aqua_secret           = "secret"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `aqua_platform_url` - (Required) The URL of the Aqua Platform.

* `aqua_key` - (Required) The API key for the Aqua Platform.

* `aqua_secret` - (Required) The API secret for the Aqua Platform.

* `aqua_auth_url` - (Optional) The URL used for authentication. Defaults to `https://api.cloudsploit.com`.

* `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)
- [Aqua Security Extension](https://marketplace.visualstudio.com/items?itemName=AquaSecurityOfficial.trivy-official)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Aqua Platform Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Aqua Platform Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Aqua Platform Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Aqua Platform Service Endpoint.

## Import

Azure DevOps Aqua Platform Service Endpoint can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_aquaplatform.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
