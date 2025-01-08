---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_jfrog_xray_v2"
description: |-
  Manages a JFrog XRay V2 service endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_jfrog_xray_v2

Manages an JFrog XRay V2 service endpoint within an Azure DevOps organization. 

~> **Note:** Using this service endpoint requires you to first install [JFrog Extension](https://marketplace.visualstudio.com/items?itemName=JFrog.jfrog-azure-devops-extension).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_jfrog_xray_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Artifactory"
  description           = "Managed by Terraform"
  url                   = "https://artifactory.my.com"
  authentication_token {
    token = "0000000000000000000000000000000000000000"
  }
}
```
Alternatively a username and password may be used.

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_jfrog_xray_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Artifactory"
  description           = "Managed by Terraform"
  url                   = "https://artifactory.my.com"
  authentication_basic {
    username = "username"
    password = "password"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `url` - (Required) URL of the Artifactory server to connect with.

~> **NOTE:** URL should not end in a slash character.

* `authentication_token` - (Optional) A `authentication_token` block as documented below.
* `authentication_basic` - (Optional) A `authentication_basic` block as documented below.
* `description` - (Optional) The Service Endpoint description.

---

A `authentication_token` block supports the following:

* `token` - Authentication Token generated through Artifactory.

---

A `authentication_basic` block supports the following:

* `username` - Artifactory Username.
* `password` - Artifactory Password.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
* [Artifactory User Token](https://docs.artifactory.org/latest/user-guide/user-token/)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the JFrog XRay V2 Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the JFrog XRay V2 Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the JFrog XRay V2 Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the JFrog XRay V2 Service Endpoint.

## Import

Azure DevOps JFrog Platform V2 Service Endpoint can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_jfrog_xray_v2.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
