---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_snyk"
description: |-
  Manages a Snyk Security Scan service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_snyk

Manages a Snyk Security Scan service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Snyk Security Scan](https://marketplace.visualstudio.com/items?itemName=Snyk.snyk-security-scan)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_snyk" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://snyk.io/"
  api_token             = "00000000-0000-0000-0000-000000000000"
  service_endpoint_name = "Example Snyk"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `server_url` - (Required) The server URL of the Snyk Security Scan.

* `api_token` - (Required) The API token of the Snyk Security Scan.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Snyk Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Snyk Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Snyk Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Snyk Service Endpoint.

## Import

Azure DevOps Snyk Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_snyk.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
