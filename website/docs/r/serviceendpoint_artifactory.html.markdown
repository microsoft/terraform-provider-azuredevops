---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_artifactory"
description: |-
  Manages an JFrog Artifactory server endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_artifactory

Manages an JFrog Artifactory server endpoint within an Azure DevOps organization. Using this service endpoint requires you to first install [JFrog Artifactory Extension](https://marketplace.visualstudio.com/items?itemName=JFrog.jfrog-artifactory-vsts-extension).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_artifactory" "example" {
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

resource "azuredevops_serviceendpoint_artifactory" "example" {
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

   _**Note: URL should not end in a slash character.**_
---

* `description` - (Optional) The Service Endpoint description.

---

* `authentication_token` - (Optional) A `authentication_token` block as defined below.

* `authentication_token` - (Optional) A `authentication_basic` block as defined below.

---

A `authentication_token` block supports the following:

* `token` - (Required) Authentication Token generated through Artifactory.

---

A `authentication_basic` block supports the following:
 
* `username` - (Required) The Username of the Artifactory.

* `password` - (Required) The Password of the Artifactory.

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

* `create` - (Defaults to 2 minutes) Used when creating the JFrog Artifactory Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the JFrog Artifactory Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the JFrog Artifactory Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the JFrog Artifactory Service Endpoint.

## Import

Azure DevOps JFrog Artifactory Service Endpoint can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_artifactory.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
