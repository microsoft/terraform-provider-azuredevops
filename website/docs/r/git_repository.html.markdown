---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository"
description: |-
  Manages a git repository within Azure DevOps organization.
---

# azuredevops_git_repository

Manages a git repository within Azure DevOps.


~>**NOTE** Importing an existing repository and running `terraform plan` will detect a difference on the `initialization` block. The `plan` and `apply` will then attempt to update the repository based on the `initialization` configurations. It may be necessary to ignore the `initialization` block from `plan` and `apply` to support configuring existing repositories imported into Terraform state.<br>

~>**NOTE**   1. `initialization.init_type` is `Uninitialized`: Changing `source_type` or `source_url` will not recreate the repository, but initialize the repository.   <br>2. `initialization.init_type` is not `Uninitialized`:
<br>&nbsp;&nbsp;&nbsp;&nbsp;1) Updating `init_type` will recreate the repository
<br>&nbsp;&nbsp;&nbsp;&nbsp;2) Updating `source_type` or `source_url` will recreate the repository


```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "Example Git Repository"
  default_branch = "refs/heads/main"
  initialization {
    init_type = "Clean"
  }
  lifecycle {
    ignore_changes = [
      # Ignore changes to initialization to support importing existing repositories
      # Given that a repo now exists, either imported into terraform state or created by terraform,
      # we don't care for the configuration of initialization against the existing resource
      initialization,
    ]
  }
}
```

## Example Usage

### Create Git repository

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}
```

### Configure existing Git repository imported into Terraform state

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "Example Git Repository"
  default_branch = "refs/heads/main"
  initialization {
    init_type = "Clean"
  }
  lifecycle {
    ignore_changes = [
      # Ignore changes to initialization to support importing existing repositories
      # Given that a repo now exists, either imported into terraform state or created by terraform,
      # we don't care for the configuration of initialization against the existing resource
      initialization,
    ]
  }
}
```

### Create Fork of another Azure DevOps Git repository

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "Example Git Repository"
  default_branch = "refs/heads/main"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository" "example-fork" {
  project_id           = azuredevops_project.example.id
  name                 = "Example Fork Repository"
  parent_repository_id = azuredevops_git_repository.example.id
  initialization {
    init_type = "Clean"
  }
}
```

### Create Import from another Git repository

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "Example Git Repository"
  default_branch = "refs/heads/main"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository" "example-import" {
  project_id = azuredevops_project.example.id
  name       = "Example Import Repository"
  initialization {
    init_type   = "Import"
    source_type = "Git"
    source_url  = "https://github.com/microsoft/terraform-provider-azuredevops.git"
  }
}
```

### Import from a Private Repository

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "Example Git Repository"
  default_branch = "refs/heads/main"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_serviceendpoint_generic_git" "example-serviceendpoint" {
  project_id            = azuredevops_project.example.id
  repository_url        = "https://dev.azure.com/org/project/_git/repository"
  username              = "username"
  password              = "<password>/<PAT>"
  service_endpoint_name = "Example Generic Git"
  description           = "Managed by Terraform"
}

resource "azuredevops_git_repository" "example-import" {
  project_id = azuredevops_project.example.id
  name       = "Example Import Existing Repository"
  initialization {
    init_type             = "Import"
    source_type           = "Git"
    source_url            = "https://dev.azure.com/example-org/private-repository.git"
    service_connection_id = azuredevops_serviceendpoint_generic_git.example-serviceendpoint.id
  }
}
```

### Disable a Git repository

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  disabled   = true
  initialization {
    init_type = "Clean"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `name` - (Required) The name of the git repository.
- `parent_repository_id` - (Optional) The ID of a Git project from which a fork is to be created.
- `disabled` - (Optional) The ability to disable or enable the repository. Defaults to `false`.
- `initialization` - (Required) An `initialization` block as documented below.

`initialization` - (Required) block supports the following:

- `init_type` - (Required) The type of repository to create. Valid values: `Uninitialized`, `Clean` or `Import`.
- `source_type` - (Optional) Type of the source repository. Used if the `init_type` is `Import`. Valid values: `Git`.
- `source_url` - (Optional) The URL of the source repository. Used if the `init_type` is `Import`.
- `service_connection_id` (Optional) The id of service connection used to authenticate to a private repository for import initialization.

## Attributes Reference

In addition to all arguments above, except `initialization`, the following attributes are exported:

- `id` - The ID of the Git repository.

- `default_branch` - The ref of the default branch. Will be used as the branch name for initialized repositories.
- `is_fork` - True if the repository was created as a fork.
- `remote_url` - Git HTTPS URL of the repository
- `size` - Size in bytes.
- `ssh_url` - Git SSH URL of the repository.
- `url` - REST API URL of the repository.
- `web_url` - Web link to the repository.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Git Repositories](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/repositories?view=azure-devops-rest-7.0)

## Import

Azure DevOps Repositories can be imported using the repo name or by the repo Guid e.g.

```sh
terraform import azuredevops_git_repository.example projectName/repoName
```

or

```sh
terraform import azuredevops_git_repository.example projectName/00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required
- **Code**: Read, Create, & Manage.
