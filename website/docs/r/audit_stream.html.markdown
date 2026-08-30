---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_audit_stream"
description: |-
  Manages an audit stream within Azure DevOps organization.
---

# azuredevops_audit_stream

Manages an audit stream within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_audit_stream" "example" {
  display_name  = "Example Audit Stream"
  consumer_type = "splunk"
  status        = "enabled"

  consumer_inputs {
    key   = "url"
    value = "https://splunk.example.com"
  }

  consumer_inputs {
    key   = "token"
    value = "your-splunk-token"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The human-readable name for the audit stream.
* `consumer_type` - (Required) The type of the consumer. Possible values: `splunk`, `azureMonitorLogs`.
* `consumer_inputs` - (Required) A list of key-value pairs of consumer inputs.
    * `key` - (Required) The key of the consumer input.
    * `value` - (Required, Sensitive) The value of the consumer input.
* `status` - (Optional) The status of the stream. Possible values: `enabled`, `disabledByUser`. Defaults to `enabled`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the audit stream.
* `created_time` - The time when the stream was created.
* `updated_time` - The time when the stream was last updated.
* `status_reason` - The reason for the current status.

## Relevant Links

- [Azure DevOps Auditing Streaming](https://docs.microsoft.com/en-us/azure/devops/organizations/audit/auditing-streaming?view=azure-devops)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Audit Stream.
* `read` - (Defaults to 1 minute) Used when retrieving the Audit Stream.
* `update` - (Defaults to 2 minutes) Used when updating the Audit Stream.
* `delete` - (Defaults to 2 minutes) Used when deleting the Audit Stream.

## Import

Azure DevOps Audit Streams can be imported using the audit stream ID, e.g.

```sh
terraform import azuredevops_audit_stream.example 0
```
