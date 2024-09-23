---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_wiki"
description: |-
  Manages Wikis within Azure DevOps project.
---

# azuredevops_wiki

Manages Wikis within Azure DevOps project.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name        = "Example Project"
  description = "Managed by Terraform"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_wiki" "example" {
  name       = "Example project wiki "
  project_id = azuredevops_project.example.id
  type       = "projectWiki"
}

resource "azuredevops_wiki" "example2" {
  name          = "Example wiki in repository"
  project_id    = azuredevops_project.example.id
  repository_id = azuredevops_git_repository.example.id
  version       = "main"
  type          = "codeWiki"
  mappedpath    = "/"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

* `name` - (Required) The name of the Wiki.

* `type` -  (Required) The type of the wiki. Possible values are `codeWiki`, `projectWiki`.

* `repository_id` - (Optional) The ID of the repository.

* `version` - (Optional) Version of the wiki.

* `mappedpath` - (Optional) Folder path inside repository which is shown as Wiki.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the wiki returned after creation in Azure DevOps.
* `remote_url` - The remote web url to the wiki.
* `url` - The REST url for this wiki.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Wiki ](https://learn.microsoft.com/en-us/rest/api/azure/devops/wiki/wikis?view=azure-devops-rest-7.1)

## Import

Azure DevOps Wiki can be imported using the `id`

```shell
terraform import azuredevops_wiki.wiki 00000000-0000-0000-0000-000000000000
```
