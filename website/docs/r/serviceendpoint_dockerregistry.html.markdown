---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_dockerregistry"
description: |-
  Manages a Docker Registry service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_dockerregistry

Manages a Docker Registry service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

# dockerhub registry service connection
resource "azuredevops_serviceendpoint_dockerregistry" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Docker Hub"
  docker_username       = "example"
  docker_email          = "email@example.com"
  docker_password       = "12345"
  registry_type         = "DockerHub"
}

# other docker registry service connection
resource "azuredevops_serviceendpoint_dockerregistry" "example-other" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Docker Registry"
  docker_registry       = "https://sample.azurecr.io/v1"
  docker_username       = "sample"
  docker_password       = "12345"
  registry_type         = "Others"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The name you will use to refer to this service connection in task inputs.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.
- `docker_registry` - (Optional) The URL of the Docker registry. (Default: "https://index.docker.io/v1/")
- `docker_username` - (Optional) The identifier of the Docker account user.
- `docker_email` - (Optional) The email for Docker account user.
- `docker_password` - (Optional) The password for the account user identified above.
- `registry_type` - (Optional) Can be "DockerHub" or "Others" (Default "DockerHub")

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)
- [Docker Registry Service Connection](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml#sep-docreg)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Docker Registry Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Docker Registry Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Docker Registry Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Docker Registry Service Endpoint.

## Import

Azure DevOps Docker Registry Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_dockerregistry.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
