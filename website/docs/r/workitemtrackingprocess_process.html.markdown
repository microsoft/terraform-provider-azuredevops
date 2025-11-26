---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_process"
description: |-
  Manages a process.
---

# azuredevops_workitemtrackingprocess_process

Manages a process.

## Example Usage

```hcl
resource "azuredevops_workitemtrackingprocess_process" "custom_agile" {
  name = "custom_agile"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc" // Agile
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required)  Name of the process.

* `parent_process_type_id` - (Required)  ID of the parent process. Changing this forces a new process to be created.

---

* `description` - (Optional)  Description of the process. Default: ""

* `is_default` - (Optional)  Is the process default? Default: false

* `is_enabled` - (Optional)  Is the process enabled? Default: true

* `reference_name` - (Optional)  Reference name of process being created. If not specified, server will assign a unique reference name. Changing this forces a new process to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the process.

* `customization_type` -  Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Processes](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/processes?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the process.
* `read` - (Defaults to 5 minutes) Used when retrieving the process.
* `update` - (Defaults to 10 minutes) Used when updating the process.
* `delete` - (Defaults to 10 minutes) Used when deleting the process.

## Import

A process can be imported using the process id, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_process.example 00000000-0000-0000-0000-000000000000
```
