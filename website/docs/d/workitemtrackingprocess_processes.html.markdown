---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_workitemtrackingprocess_processes"
description: |-
  Gets information about existing processes.
---

# Data Source: azuredevops_workitemtrackingprocess_processes

Use this data source to access information about existing processes.

## Example Usage

```hcl
data "azuredevops_workitemtrackingprocess_processes" "all" {

}

output "id" {
  value = data.azuredevops_workitemtrackingprocess_processes.all.id
}
```

## Arguments Reference

The following arguments are supported:

* `expand` - (Optional)  Specifies the expand option when getting the processes. Default: "none"

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the process.

* `processes` - A `processes` block as defined below. A list of all processes including system and inherited.

---

A `processes` block exports the following:

* `customization_type` -  Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before.

* `description` -  Description of the process.

* `id` -  The ID of the process.

* `is_default` -  Is the process default?

* `is_enabled` -  Is the process enabled?

* `name` -  Name of the process.

* `parent_process_type_id` -  ID of the parent process.

* `projects` - A `projects` block as defined below. Returns associated projects when using the 'projects' expand option.

* `reference_name` -  Reference name of process being created. If not specified, server will assign a unique reference name.

---

A `projects` block exports the following:

* `description` -  Description of the project.

* `id` -  The ID of the project.

* `name` -  Name of the project.

* `url` -  Url of the project.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Processes - List](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/processes/list?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the process.
