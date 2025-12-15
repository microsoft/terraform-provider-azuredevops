---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_workitemtrackingprocess_workitemtype"
description: |-
  Gets information about an existing work item type for a process.
---

# Data Source: azuredevops_workitemtrackingprocess_workitemtype

Use this data source to access information about an existing work item type for a process.

## Example Usage

```hcl
data "azuredevops_workitemtrackingprocess_workitemtype" "example" {
  process_id     = "f22ab9cc-acad-47ab-b31d-e43ef8d72b89"
  reference_name = "MyProcess.Example"
}

output "name" {
  value = data.azuredevops_workitemtrackingprocess_workitemtype.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required)  The ID of the process.

* `reference_name` - (Required)  The reference name of the work item type.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the work item type.

* `color` -  Color hexadecimal code to represent the work item type.

* `customization` -  Indicates the type of customization on this work item type.

* `description` -  Description of the work item type.

* `icon` -  Icon to represent the work item type.

* `parent_work_item_reference_name` - Reference name of the parent work item type.

* `is_disabled` -  Indicates if the work item type is disabled.

* `name` -  Name of the work item type.

* `url` -  URL of the work item type.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Work Item Types - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/work-item-types/get?view=azure-devops-rest-7.1&tabs=HTTP)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the work item type.
