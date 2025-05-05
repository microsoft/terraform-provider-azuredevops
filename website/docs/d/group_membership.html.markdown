---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group_membership"
description: |-
  Use this data source to access information about an existing Group Memberships within Azure DevOps
---

# azuredevops_group_membership

Use this data source to access information about an existing Group Memberships within Azure DevOps

## Example Usage

```hcl
data "azuredevops_group_membership" "example" {
  group_descriptor = "groupdescroptpr"
}
```

## Argument Reference

The following arguments are supported:

* `group_descriptor` - (Required) The descriptor of the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The descriptor of the group.

* `members` - A list of user or group descriptors.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Memberships](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Group membership.
