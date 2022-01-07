---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_auditstream_azureeventgrid"
description: |-
  Manages an Azure EventGrid audit stream within and Azure DevOps organization.
---

# azuredevops_auditstream_azureeventgrid

Manages an Azure EventGrid audit stream within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_auditstream_azureeventgrid" "example" {
  topic_url        = "https://topic1.westus2-1.eventgrid.azure.net/api/events"
  access_key       = "0000000000000000000000000000000000000000"
}
```

## Arguments Reference

The following arguments are supported:

- `topic_url` - (Required) Url of your Azure Event Grid topic that will send events to. It should look like `https://topic1.westus2-1.eventgrid.azure.net/api/events`.
- `access_key` - (Required) Access key found in the settings of the Azure Event Grid topic.
- `days_to_backfill` - (Optional) The number of days of previously recorded audit data that will be replayed into the stream. A value of zero will result in only new events being streamed. Defaults to `0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the audit stream.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Audit Streams](https://docs.microsoft.com/en-us/rest/api/azure/devops/audit/?view=azure-devops-rest-6.0)

## Import

Azure DevOps Audit Streams can be imported using the audit stream ID , e.g.

```shell
$ terraform import azuredevops_auditstream_azureeventgrid.example 10
```
