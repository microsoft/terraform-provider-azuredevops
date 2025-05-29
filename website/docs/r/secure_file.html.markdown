---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_secure_file"
description: |-
  Manages secure files within Azure DevOps project.
---

# azuredevops_secure_file

Manages secure files within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_secure_file" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "my-secure-file.txt"
  content    = file("./my-secure-file.txt")
  properties = {
    environment = "production"
  }
  allow_access = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `name` - (Required) The name of the secure file. Must be unique within the project.
* `content` - (Required) The content of the secure file. This is the actual file data that will be stored securely.
* `properties` - (Optional) Key-value map of properties to associate with the secure file. The provider automatically adds `file_hash_sha1` and `file_hash_sha256`.
* `allow_access` - (Optional) Boolean that indicate if this secure file is shared by all pipelines of this project. Defaults to `false`. If set to `true`, the secure file can be used in all pipelines without needing explicit permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the secure file returned after creation in Azure DevOps.
* `file_hash_sha1` - SHA1 hash of the file content. Computed from the content during creation.
* `file_hash_sha256` - SHA256 hash of the file content. Computed from the content during creation.

## Relevant Links

API documentation for secure files is not available in the Azure DevOps REST API documentation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the secure file.
* `read` - (Defaults to 5 minute) Used when retrieving the secure file.
* `update` - (Defaults to 10 minutes) Used when updating the secure file.
* `delete` - (Defaults to 10 minutes) Used when deleting the secure file.

## Import

Azure DevOps secure files can be imported using the project name/secure file ID or by the project Guid/secure file ID, e.g.

```sh
terraform import azuredevops_secure_file.example "Example Project/00000000-0000-0000-0000-000000000000"
```

or

```sh
terraform import azuredevops_secure_file.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

_Note that the content of secure files is not imported and will not be present in the state._

## PAT Permissions Required

- **Secure Files**: Read, Create, & Manage
- **Build**: Read & execute
- **Project and Team**: Read

