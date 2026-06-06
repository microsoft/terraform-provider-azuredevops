---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_audit_stream"
description: |-
  Find an existing audit stream within Azure DevOps organization.
---

# azuredevops_audit_stream

Find an existing audit stream within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_audit_stream" "example" {
  display_name = "Example Audit Stream"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The human-readable name for the audit stream.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the audit stream.
* `consumer_type` - The type of the consumer.
* `status` - The status of the stream.
* `consumer_inputs` - A list of key-value pairs of consumer inputs.
    * `key` - The key of the consumer input.
    * `value` - The value of the consumer input.
* `created_time` - The time when the stream was created.
* `updated_time` - The time when the stream was last updated.
* `status_reason` - The reason for the current status.
