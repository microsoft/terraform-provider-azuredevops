---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_audit_streams"
description: |-
  Lists all audit streams within Azure DevOps organization.
---

# azuredevops_audit_streams

Lists all audit streams within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_audit_streams" "example" {
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the data source.
* `streams` - A list of all configured audit streams in the organization. Each stream block contains:
    * `id` - The unique ID of the audit stream.
    * `display_name` - The human-readable name for the audit stream.
    * `consumer_type` - The type of the consumer.
    * `status` - The status of the stream.
    * `consumer_inputs` - A list of key-value pairs of consumer inputs.
        * `key` - The key of the consumer input.
        * `value` - The value of the consumer input.
    * `created_time` - The time when the stream was created.
    * `updated_time` - The time when the stream was last updated.
    * `status_reason` - The reason for the current status.
