---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_externaltfs"
description: |-
  Manages an Azure Repository/Team Foundation Server service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_externaltfs

Manages an Azure Repository/Team Foundation Server service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_externaltfs" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example External TFS Name"
  connection_url        = "https://dev.azure.com/myorganization"
  description           = "Managed by Terraform"

  auth_personal {
    # Also can be set with AZDO_PERSONAL_ACCESS_TOKEN environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `connection_url` - (Required) Azure DevOps Organization or TFS Project Collection Url.

* `auth_personal` - (Required) An `auth_personal` block as documented below. Allows connecting using a personal access token.

---

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

---

A `auth_personal` block supports the following:

* `personal_access_token` - (Required) The Personal Access Token for Azure DevOps Organization.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Azure Repository/Team Foundation Server Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Azure Repository/Team Foundation Server Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Azure Repository/Team Foundation Server Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Azure Repository/Team Foundation Server Service Endpoint.

## Import

Azure DevOps Azure Repository/Team Foundation Server Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_externaltfs.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
