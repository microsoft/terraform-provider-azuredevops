---
layout: "azuredevops"
page_title: "AzureDevops: data.azuredevops_secure_file"
description: |-
  Retrieves information about a Secure File in an Azure DevOps project.
---

# data "azuredevops_secure_file"

Retrieves information about a Secure File in an Azure DevOps project.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_secure_file" "example" {
  project_id = azuredevops_project.example.id
  name       = "my-secure-file.txt"
  content    = filebase64("./my-secure-file.txt")
}

data "azuredevops_secure_file" "example" {
  project_id = azuredevops_project.example.id
  name       = azuredevops_secure_file.example.name
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `name` - (Required) The name of the secure file to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the secure file.
* `file_hash_sha1` - SHA1 hash of the file content.
* `file_hash_sha256` - SHA256 hash of the file content.
* `properties` - Key-value map of properties associated with the secure file.

## Relevant Links

API is not documented in the Azure DevOps REST API documentation.

## PAT Permissions Required

- **Secure Files**: Read
- **Project and Team**: Read

