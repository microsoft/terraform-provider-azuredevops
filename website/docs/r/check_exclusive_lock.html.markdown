---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_exclusive_lock"
description: |-
  Manages a Exclusive Lock Check.
---

# azuredevops_check_exclusive_lock

Manages a Exclusive Lock Check.

Adding an exclusive lock will only allow a single stage to utilize this resource at a time. If multiple stages are waiting on the lock, only the latest will run. All others will be canceled.

## Example Usage

### Add Exclusive Lock to an environment

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_serviceendpoint_generic" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://some-server.example.com"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Example Generic"
  description           = "Managed by Terraform"
}

data "azuredevops_users" "example" {
  principal_name = "someone@somewhere.com"
}

resource "azuredevops_check_exclusive_lock" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_serviceendpoint_generic.example.id
  target_resource_type = "endpoint"

  timeout = 43200
}
```

### Protect an environment

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_environment" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Environment"
}

resource "azuredevops_group" "example" {
  display_name = "some-azdo-group"
}

resource "azuredevops_check_exclusive_lock" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_environment.example.id
  target_resource_type = "environment"

  timeout = 43200
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The project ID. Changing this forces a new Exclusive Lock Check to be created.

* `target_resource_id` - (Required) The ID of the resource being protected by the check. Changing this forces a new Exclusive Lock to be created.

* `target_resource_type` - (Required) The type of resource being protected by the check. Valid values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`. Changing this forces a new Exclusive Lock to be created.
 
---

* `timeout` - (Optional) The timeout in minutes for the exclusive lock.  Defaults to `43200`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Exclusive Lock.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Exclusive Lock.
* `read` - (Defaults to 1 minute) Used when retrieving the Exclusive Lock.
* `update` - (Defaults to 2 minutes) Used when updating the Exclusive Lock.
* `delete` - (Defaults to 2 minutes) Used when deleting the Exclusive Lock.

## Import

Importing this resource is not supported.
