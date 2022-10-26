---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_sonarcloud"
description: |-
  Manages the SonarCloud service endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_sonarcloud
Manages a SonarCloud service endpoint within Azure DevOps. 

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_sonarcloud" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example SonarCloud"
  token                 = "0000000000000000000000000000000000000000"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `token` - (Required) Authentication Token generated through SonarCloud (go to My Account > Security > Generate Tokens).
* `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)
- [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
- [SonarCloud User Token](https://docs.sonarcloud.org/latest/user-guide/user-token/)

## Import
Azure DevOps Service Endpoint SonarCloud can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_sonarcloud.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
