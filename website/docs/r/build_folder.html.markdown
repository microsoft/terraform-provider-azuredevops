---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_build_folder"
description: |-
  Manages a Build Folder.
---

# azuredevops_build_folder

Manages a Build Folder.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_build_folder" "example" {
  project_id  = azuredevops_project.example.id
  path        = "\\ExampleFolder"
  description = "ExampleFolder description"
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project in which the folder will be created.

* `path` - (Required) The folder path.

---

* `description` - (Optional) Folder Description.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Build Folder.
* `read` - (Defaults to 5 minute) Used when retrieving the Build Folder.
* `update` - (Defaults to 30 minutes) Used when updating the Build Folder.
* `delete` - (Defaults to 30 minutes) Used when deleting the Build Folder.

## Import

Build Folders can be imported using the `project name/path` or `project id/path`, e.g.

```shell
terraform import azuredevops_build_folder.example "Example Project/\\ExampleFolder"
```

or

```shell
terraform import azuredevops_build_folder.example 00000000-0000-0000-0000-000000000000/\\ExampleFolder
```
