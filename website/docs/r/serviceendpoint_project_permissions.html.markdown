---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_project_permissions"
description: |-
  Manages a Service Endpoint sharing with projects.
---

# azuredevops_serviceendpoint_project_permissions

Manages a Service Endpoint sharing with projects.

## Example Usage

```hcl
resource "azuredevops_project" "example1" {
  name               = "Example Project 1"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
  features = {
    testplans = "disabled"
    artifacts = "disabled"
  }
}

resource "azuredevops_project" "example2" {
  name               = "Example Project 2"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
  features = {
    testplans = "disabled"
    artifacts = "disabled"
  }
}

resource "azuredevops_project" "example3" {
  name               = "Example Project 3"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
  features = {
    testplans = "disabled"
    artifacts = "disabled"
  }
}

resource "azuredevops_serviceendpoint_azuredevops" "example" {
  project_id            = azuredevops_project.example1.id
  service_endpoint_name = "Example Azure DevOps"
  org_url               = "https://dev.azure.com/testorganization"
  release_api_url       = "https://vsrm.dev.azure.com/testorganization"
  personal_access_token = "0000000000000000000000000000000000000000000000000000"
  description           = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_project_permissions" "example-share" {
  project_id = azuredevops_project.example1.id
  serviceendpoint_id = azuredevops_serviceendpoint_azuredevops.example.id

  project_reference {
    project_id            = azuredevops_project.example2.id
    service_endpoint_name = "service-connection-shared"
    description           = "Service Connection Shared by Terraform - Cluster Two"
  }

  project_reference {
    project_id            = azuredevops_project.example3.id
    service_endpoint_name = "service-connection-shared"
    description           = "Service Connection Shared by Terraform - Cluster Three"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Endpoint sharing with projects to be created.

* `service_endpoint_id` - (Required) The ID of the endpoint. Changing this forces a new Service Endpoint sharing with projects to be created.

---

* `project_reference` - (Optional) One or more `project_reference` blocks as defined below.

---

A `project_reference` block supports the following:

* `project_id` - (Required) The ID of the project.

* `description` - (Optional) Description of the service endpoint..

* `service_endpoint_name` - (Optional) Name of the service endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Endpoint sharing with projects.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Service Endpoint sharing with projects.
* `read` - (Defaults to 1 minute) Used when retrieving the Service Endpoint sharing with projects.
* `update` - (Defaults to 2 minutes) Used when updating the Service Endpoint sharing with projects.
* `delete` - (Defaults to 2 minutes) Used when deleting the Service Endpoint sharing with projects.

## Import

Service Endpoint sharing with projectss can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_project_permissions.example 00000000-0000-0000-0000-000000000000
```
