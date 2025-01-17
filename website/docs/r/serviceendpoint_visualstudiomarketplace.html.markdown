---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_visualstudiomarketplace"
description: |-
  Manages a Visual Studio Marketplace service endpoint within Azure DevOps organization. Packaging and publishing Azure Devops and Visual Studio extensions to the Visual Studio Marketplace.
---

# azuredevops_serviceendpoint_visualstudiomarketplace

Manages a Visual Studio Marketplace service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Azure DevOps Extension Tasks](https://marketplace.visualstudio.com/items?itemName=ms-devlabs.vsts-developer-tools-build-tasks)

## Example Usage

###  Authorize with token
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Marketplace"
  url                   = "https://markpetplace.com"
  authentication_token {
    token = "token"
  }
  description = "Managed by Terraform"
}
```

### Authorize with username and password

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_visualstudiomarketplace" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Marketplace"
  url                   = "https://markpetplace.com"
  authentication_basic {
    username = "username"
    password = "password"
  }
  description = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `url` - (Required) The server URL for Visual Studio Marketplace.

---

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

* `authentication_token` - (Optional) An `authentication_token` block as documented below.

* `authentication_basic` - (Optional) An `authentication_basic` block as documented below.

~> **NOTE:** `authentication_basic` and `authentication_token` conflict with each other, only one is required.

---

An `authentication_token` block supports the following:

* `token` - The Personal Access Token.

---

An `authentication_basic` block supports the following:

* `username` - The username of the marketplace.

* `password` - The password of the marketplace.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.

* `project_id` - The ID of the project.

* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Visual Studio Marketplace Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Visual Studio Marketplace Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Visual Studio Marketplace Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Visual Studio Marketplace Service Endpoint.

## Import

Azure DevOps Visual Studio Marketplace Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_visualstudiomarketplace.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
