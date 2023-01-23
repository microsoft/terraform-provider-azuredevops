---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_jfrog_artifactory_v2"
description: |-
  Manages a JFrog Artifactory V2 server endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_jfrog_artifactory_v2
Manages a JFrog Artifactory V2 server endpoint within an Azure DevOps organization.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_jfrog_artifactory_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example JFrog Artifactory V2"
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

resource "azuredevops_serviceendpoint_jfrog_artifactory_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example JFrog Artifactory V2"
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

   _Note: URL should not end in a slash character._
* either `authentication_token` or `authentication_basic` (one is required)
  * `authentication_token`
    * `token` - Authentication Token generated through Artifactory.
  * `authentication_basic`
      * `username` - Artifactory Username.
      * `password` - Artifactory Password.
* `description` - (Optional) The Service Endpoint description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
* [Artifactory User Token](https://docs.artifactory.org/latest/user-guide/user-token/)

## Import
Azure DevOps Service Endpoint JFrog Artifactory V2 can be imported using the **projectID/serviceEndpointID**, e.g.

```sh
terraform import azuredevops_serviceendpoint_jfrog_artifactory_v2.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
