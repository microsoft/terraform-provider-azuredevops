---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_gcp_terraform"
description: |-
  Manages a GCP service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_gcp_terraform
Manages a GCP service endpoint for Terraform within Azure DevOps. Using this service endpoint requires you to first install the [Azure Terraform Extension for Azure DevOps from Microsoft DevLabs](https://marketplace.visualstudio.com/items?itemName=ms-devlabs.custom-terraform-tasks).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_gcp_terraform" "example" {
  project_id            = azuredevops_project.example.id
  token_uri             = "https://oauth2.example.com/token"
  client_email          = "gcp-sa-example@example.iam.gserviceaccount.com"
  private_key           = "0000000000000000000000000000000000000"
  service_endpoint_name = "Example GCP Terraform extension"
  gcp_project_id        = "Example GCP Project"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `private_key` - (Required) The client email field in the JSON key file for creating the JSON Web Token.

* `token_uri` - (Required) The token uri field in the JSON key file for creating the JSON Web Token.

* `gcp_project_id` - (Required) GCP project associated with the Service Connection.

---

* `client_email` - (Optional) The client email field in the JSON key file for creating the JSON Web Token.
 
* `scope` - (Optional) Scope to be provided.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
- [Azure DevOps Service REST API 7.1 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.1)

## Import
Azure DevOps Service Endpoint GCP can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
 terraform import azuredevops_serviceendpoint_gcp_terraform.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
