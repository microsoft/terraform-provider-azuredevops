---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_securityrole_assignment"
description: |-
  Manages assignment of security roles to various resources within Azure DevOps organization.
---

# azuredevops_securityrole_assignment

Manages assignment of security roles to various resources within Azure DevOps organization.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_environment" "example" {
  project_id  = azuredevops_project.example.id
  name        = "Example Environment"
  description = "Example pipeline deployment environment"
}

resource "azuredevops_group" "example" {
  scope        = azuredevops_project.example.id
  display_name = "Example group"
  description  = "Description of example group"
}

resource "azuredevops_securityrole_assignment" "example" {
  scope       = "distributedtask.environmentreferencerole"
  resource_id = format("%s_%s", azuredevops_project.example.id, azuredevops_environment.example.id)
  identity_id = azuredevops_group.example.origin_id
  role_name   = "Administrator"
}
```

## Argument Reference

The following arguments are supported:

* `scope` - (Required) The scope in which this assignment should exist.

* `resource_id` - (Required) The ID of the resource on which the role is to be assigned. Changing this forces a new resource to be created.

* `identity_id` - (Required) The ID of the identity to authorize.

* `role_name` - (Required) Name of the role to assign.

## Attributes Reference

No attributes are exported

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Security Role Assignment.
* `read` - (Defaults to 5 minute) Used when retrieving the Security Role Assignment.
* `update` - (Defaults to 10 minutes) Used when updating the Security Role Assignment.
* `delete` - (Defaults to 10 minutes) Used when deleting the Security Role Assignment.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Authorize Definition Resource](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/resources/authorize%20definition%20resources?view=azure-devops-rest-7.0)
