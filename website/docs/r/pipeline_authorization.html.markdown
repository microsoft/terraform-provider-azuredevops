---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_pipeline_authorization"
description: |-
  Manages Pipeline Authorizations within Azure DevOps Project.
---

# azuredevops_pipeline_authorization

Manage pipeline access permissions to resources.

~> **Note** This resource is a replacement for `azuredevops_resource_authorization`.  Pipeline authorizations managed by `azuredevops_resource_authorization` can also be managed by this resource.

~> **Note** If both "All Pipeline Authorization" and "Custom Pipeline Authorization" are configured, "All Pipeline Authorization" has higher priority.


## Example Usage

### Authorization for all pipelines

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "example" {
  name           = "Example Pool"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "example" {
  project_id    = azuredevops_project.example.id
  agent_pool_id = azuredevops_agent_pool.example.id
}

resource "azuredevops_pipeline_authorization" "example" {
  project_id  = azuredevops_project.example.id
  resource_id = azuredevops_agent_queue.example.id
  type        = "queue"
}
```

### Authorization for specific pipeline

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "example" {
  name           = "Example Pool"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "example" {
  project_id    = azuredevops_project.example.id
  agent_pool_id = azuredevops_agent_pool.example.id
}

data "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Project"
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Pipeline"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.example.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "example" {
  project_id  = azuredevops_project.example.id
  resource_id = azuredevops_agent_queue.example.id
  type        = "queue"
  pipeline_id = azuredevops_build_definition.example.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The  ID of the project. Changing this forces a new resource to be created

* `resource_id` - (Required) The ID of the resource to authorize. Changing this forces a new resource to be created

* `type` - (Required) The type of the resource to authorize. Possible values are: `endpoint`, `queue`, `variablegroup`, `environment`, `repository`. Changing this forces a new resource to be created

~> **Note** `repository` is for AzureDevOps repository. To authorize repository other than 
    Azure DevOps like GitHub you need to use service connection(`endpoint`)  to connect and authorize.      
    Typical process for connecting to GitHub:    
    **Pipeline  <-----> Service Connection(`endpoint`) <-----> GitHub Repository**

---

* `pipeline_id` - (Optional) The ID of the pipeline. If not configured, all pipelines will be authorized. Changing this forces a new resource to be created.

* `pipeline_project_id` - (Optional) The ID of the project where the pipeline exists. Defaults to `project_id` if not specified. Changing this forces a new resource to be created

## Attributes Reference

No attributes are exported

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Pipeline Authorization.
* `read` - (Defaults to 1 minute) Used when retrieving the Pipeline Authorization.
* `update` - (Defaults to 2 minutes) Used when updating the Pipeline Authorization.
* `delete` - (Defaults to 2 minutes) Used when deleting the Pipeline Authorization.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Pipeline Permissions](https://learn.microsoft.com/en-us/rest/api/azure/devops/approvalsandchecks/pipeline-permissions?view=azure-devops-rest-7.1)
 
