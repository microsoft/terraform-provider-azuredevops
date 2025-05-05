---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_npm"
description: |-
  Manages a NPM service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_npm

Manages a NPM service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_npm" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example npm"
  url                   = "https://registry.npmjs.org"
  access_token          = "00000000-0000-0000-0000-000000000000"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `url` - (Required) URL of the npm registry to connect with.

* `access_token` - (Required) The access token for npm registry.

---

* `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)
- [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
- [npm User Token](https://docs.npmjs.com/about-access-tokens)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the NPM Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the NPM Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the NPM Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the NPM Service Endpoint.

## Import

Azure DevOps NPM Service Endpoint can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_npm.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
