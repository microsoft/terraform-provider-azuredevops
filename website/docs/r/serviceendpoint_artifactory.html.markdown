---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_artifactory"
description: |-
  Manages an Artifactory server endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_artifactory
Manages an Artifactory server endpoint within an Azure DevOps organization. 

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_artifactory" "serviceendpoint" {

  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample Artifactory"
  description           = "Managed by Terraform"
  url                   = "https://artifactory.my.com"
  authentication_token {
      token      = "0000000000000000000000000000000000000000"
  }
}
```
Alternatively a username and password may be used.

```hcl
resource "azuredevops_serviceendpoint_artifactory" "serviceendpoint" {

  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample Artifactory"
  description           = "Managed by Terraform"
  url                   = "https://artifactory.my.com"
  authentication_basic {
      username              = "sampleuser"
      password              = "0000000000000000000000000000000000000000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
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
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
* [Artifactory User Token](https://docs.artifactory.org/latest/user-guide/user-token/)

## Import
Azure DevOps Service Endpoint Artifactory can be imported using the **projectID/serviceEndpointID**, e.g.

```shell
terraform import azuredevops_serviceendpoint_artifactory.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
