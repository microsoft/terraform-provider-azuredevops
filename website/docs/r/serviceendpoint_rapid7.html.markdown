---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_rapid7"
description: |-
  Manages a Rapid7 Insight App sec service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_rapid7
Manages a Rapid7 service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_rapid7" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Rapid 7 Test"
  description           = "Managed by AzureDevOps"

  region                = "eu"
  auth_token            = "000000000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `region` - (Required) Rapid7 app region.
* `auth_token` - (Required) Rapid7 Authentication token.
* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Rapid7 Azure Devops Extension](https://github.com/rapid7/insightappsec-azure-devops-extension)
* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)


```
```
