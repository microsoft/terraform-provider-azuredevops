---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_auditstream_azuremonitorlogs"
description: |-
  Manages an Azure Monitor Logs audit stream within and Azure DevOps organization.
---

# azuredevops_auditstream_azuremonitorlogs

Manages an Azure Monitor Logs audit stream within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_auditstream_azuremonitorlogs" "example" {
  workspace_id     = "00000000-0000-0000-0000-000000000000"
  shared_key       = "0000000000000000000000000000000000000000"
}
```

## Arguments Reference

The following arguments are supported:

- `workspace_id` - (Required) Workspace Id of the Azure Monitor Logs instance. It should look like `00000000-0000-0000-0000-000000000000`.
- `shared_key` - (Required) Shared Key to authenticate to the Azure Monitor Logs instance.
- `days_to_backfill` - (Optional) The number of days of previously recorded audit data that will be replayed into the stream. A value of zero will result in only new events being streamed. Defaults to `0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the audit stream.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Audit Streams](https://docs.microsoft.com/en-us/rest/api/azure/devops/audit/?view=azure-devops-rest-6.0)

## Import

Azure DevOps Audit Streams can be imported using the audit stream ID , e.g.

```shell
terraform import azuredevops_auditstream_azuremonitorlogs.example 10
```
