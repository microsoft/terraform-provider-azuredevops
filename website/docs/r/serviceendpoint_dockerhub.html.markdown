---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_dockerhub"
description: |-
  Manages a Docker Hub service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_dockerhub
Manages a Docker Hub service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_dockerhub" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "Sample Docker Hub"

    docker_username        = "sample"
    docker_email           = "email@example.com" 
    docker_password        = "12345"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `description` - (Required) The Service Endpoint description.
* `docker_username` - (Required) The username for Docker Hub account.
* `docker_email` - (Required) The email for Docker Hub account.
* `docker_password` - (Required) The password for Docker Hub account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
