---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group"
description: |-
  Use this data source to access information about an existing Group within Azure DevOps.
---

# Data Source: azuredevops_group
Use this data source to access information about an existing Group within Azure DevOps

## Example Usage

```hcl
data "azuredevops_group" "test" {
    project_id = azuredevops_project.project.id
    name       = "Test Group"
}

output "group_id" {
    value = "${data.azuredevops_group.test.id}"
}
output "group_descriptor" {
    value = "${data.azuredevops_group.test.descriptor}"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Project Id.
* `name` - (Required) The Group Name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID for this resource is the group descriptor. See below.
* `descriptor` - The Descriptor is the primary way to reference the graph subject. This field will uniquely identify the same graph subject across both Accounts and Organizations.

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Groups - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups/get?view=azure-devops-rest-5.1)