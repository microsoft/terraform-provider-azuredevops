---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_wiki_page"
description: |-
  Manages Wiki pages within Azure DevOps project.
---

# azuredevops_wiki_page

Manages Wiki pages within Azure DevOps project.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name        = "Example Project"
  description = "Managed by Terraform"
}

resource "azuredevops_wiki" "example" {
  name       = "Example project wiki "
  project_id = azuredevops_project.example.id
  type       = "projectWiki"
}

resource "azuredevops_wiki_page" "example" {
  project_id = azuredevops_project.example.id
  wiki_id = azuredevops_wiki.example.id
  path = "/page"
  content = "content"
}

```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

* `wiki_id` - (Required) The ID of the Wiki.

* `path` -  (Required) The path of the wiki page.

* `content` - (Required) The content of the wiki page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the wiki page returned after creation in Azure DevOps.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Wiki Page](https://learn.microsoft.com/en-us/rest/api/azure/devops/wiki/pages?view=azure-devops-rest-7.1)
