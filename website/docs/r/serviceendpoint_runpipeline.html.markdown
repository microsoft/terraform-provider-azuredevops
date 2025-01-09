---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_runpipeline"
description: |-
  Manages a Azure DevOps plugin RunPipeline.
---

# azuredevops_serviceendpoint_runpipeline

Manages a Azure DevOps Run Pipeline Service Connection service endpoint within Azure DevOps. Allows to run downstream pipelines, monitoring their execution, collecting and consolidating artefacts produced in the delegate pipelines (yaml block `task: RunPipelines@1`). More details on Marketplace page: [RunPipelines](https://marketplace.visualstudio.com/items?itemName=CSE-DevOps.RunPipelines)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_runpipeline" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Pipeline Runner"
  organization_name     = "Organization Name"
  auth_personal {
    # Also can be set with AZDO_PERSONAL_ACCESS_TOKEN environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  description = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `organization_name` - (Required) The organization name used for `Organization Url` and `Release API Url` fields.
- `auth_personal` - (Required) An `auth_personal` block as documented below. Allows connecting using a personal access token.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

`auth_personal` block supports the following:

- `personal_access_token` - (Required) The Personal Access Token for Azure DevOps Pipeline. It also can be set with AZDO_PERSONAL_ACCESS_TOKEN environment variable.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Run Pipeline Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Run Pipeline Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Run Pipeline Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Run Pipeline Service Endpoint.

## Import

Azure DevOps Run Pipeline Service Endpoint can be imported using the `project id`, `service connection id`, e.g.

```sh
terraform import azuredevops_serviceendpoint_runpipeline.example projectID/00000000-0000-0000-0000-000000000000
```
