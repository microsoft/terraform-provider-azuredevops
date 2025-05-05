---
layout: "azuredevops"
page_title: "AzureDevops: extension"
description: |-
  Manages Extension within Azure DevOps organization.
---

# azuredevops_extension

Manages extension within Azure DevOps organization.

## Example Usage

### Install Extension
```hcl
resource "azuredevops_extension" "example" {
  extension_id = "extension ID"
  publisher_id = "publisher ID"
}
```

## Argument Reference

The following arguments are supported:

* `extension_id` - (Required) The publisher ID of the extension.

* `publisher_id` - (Required) The extension ID of the extension.

---

* `disabled` - (Optional) Whether to disable the extension.

* `version`- (Optional) The version of the extension.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Extension.

* `version` - The version of the Extension.

* `extension_name` - The name of the extension.

* `publisher_name` - The name of the publisher.

* `scope` - List of all oauth scopes required by this extension.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Extension Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/extensionmanagement/installed-extensions?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the extension.
* `read` - (Defaults to 2 minute) Used when retrieving the extension.
* `update` - (Defaults to 5 minutes) Used when updating the extension.
* `delete` - (Defaults to 5 minutes) Used when deleting the extension.

## Import

Azure DevOps Extension can be imported using the publisher ID and extension ID:

```sh
terraform import azuredevops_extension.example publisherId/extensionId
```
