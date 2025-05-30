---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_generic_git"
description: |-
  Manages an Other Git service endpoint within Azure DevOps, which can be used to authenticate to any external git service
  using basic authentication via a username and password.
---

# azuredevops_serviceendpoint_generic_git

Manages an Other Git service endpoint within Azure DevOps, which can be used to authenticate to any external git service
using basic authentication via a username and password. This is mostly useful for importing private git repositories.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_generic_git" "example" {
  project_id            = azuredevops_project.example.id
  repository_url        = "https://dev.azure.com/org/project/_git/repository"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Example Generic Git"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The name of the service endpoint.

* `repository_url` - (Required) The URL of the repository associated with the service endpoint.

---

* `username` - (Optional) The username used to authenticate to the git repository.

* `password` - (Optional) The PAT or password used to authenticate to the git repository.

    ~> **Note** For AzureDevOps Git, PAT should be used as the password.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

* `enable_pipelines_access` - (Optional) A value indicating whether or not to attempt accessing this git server from Azure Pipelines.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The name of the service endpoint.
* `enable_pipelines_access` - A value indicating whether or not to attempt accessing this git server from Azure Pipelines.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Other Git Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Other Git Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Other Git Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Other Git Service Endpoint.

## Import

Azure DevOps Other Git Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_generic_git.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
