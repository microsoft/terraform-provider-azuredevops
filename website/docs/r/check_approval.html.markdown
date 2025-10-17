---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_approval"
description: |-
  Manages an Approval Check.
---

# azuredevops_check_approval

Manages a Approval Check.

## Example Usage

### Protect a service connection

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

resource "azuredevops_check_approval" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_serviceendpoint_generic.example.id
  target_resource_type = "endpoint"

  requester_can_approve = false
  approvers = [
    one(data.azuredevops_users.example.users).id,
  ]

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

resource "azuredevops_check_approval" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_environment.example.id
  target_resource_type = "environment"

  requester_can_approve = true
  approvers = [
    azuredevops_group.example.group_id,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The project ID. Changing this forces a new Approval Check to be created.

* `target_resource_id` - (Required) The ID of the resource being protected by the check. Changing this forces a new Approval Check to be created.

* `target_resource_type` - (Required) The type of resource being protected by the check. Valid values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`. Changing this forces a new Approval Check to be created.

* `approvers` - (Required) Specifies a list of approver IDs.

---

* `instructions` - (Optional) The instructions for the approvers.

* `minimum_required_approvers` - (Optional) The minimum number of approvers. This property is applicable when there is more than 1 approver.

* `requester_can_approve` - (Optional) Can the requestor approve? Defaults to `false`.

* `timeout` - (Optional) The timeout in minutes for the approval.  Defaults to `43200`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the check.
* `version` - The version of the check.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeout) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Approval Check.
* `read` - (Defaults to 1 minute) Used when retrieving the Approval Check.
* `update` - (Defaults to 2 minutes) Used when updating the Approval Check.
* `delete` - (Defaults to 2 minutes) Used when deleting the Approval Check.

## Import

Importing this resource is not supported.
