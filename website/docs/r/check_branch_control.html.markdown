---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_branch_control"
description: |-
  Manages a branch control check.
---

# azuredevops_check_branch_control

Manages a branch control check on a resource within Azure DevOps.

## Example Usage

### Protect a service connection

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_serviceendpoint_generic" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://some-server.example.com"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Example Generic"
  description           = "Managed by Terraform"
}

resource "azuredevops_check_branch_control" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_serviceendpoint_generic.example.id
  target_resource_type = "endpoint"
  allowed_branches     = "refs/heads/main, refs/heads/features/*"

  timeout = 1440
}
```

### Protect an environment

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_environment" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Environment"
}

resource "azuredevops_check_branch_control" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_environment.example.id
  target_resource_type = "environment"
  allowed_branches     = "refs/heads/main, refs/heads/features/*"
}
```

### Protect an agent queue

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_agent_pool" "example" {
  name = "example-pool"
}

resource "azuredevops_agent_queue" "example" {
  project_id    = azuredevops_project.example.id
  agent_pool_id = azuredevops_agent_pool.example.id
}

resource "azuredevops_check_branch_control" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_agent_queue.example.id
  target_resource_type = "queue"
  allowed_branches     = "refs/heads/main, refs/heads/features/*"
}
```

### Protect a repository

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_check_branch_control" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = "${azuredevops_project.example.id}.${azuredevops_git_repository.example.id}"
  target_resource_type = "repository"
  allowed_branches     = "refs/heads/main, refs/heads/features/*"
}
```

### Protect a variable group

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Variable Group"
  description  = "Example Variable Group Description"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }

  variable {
    name         = "key2"
    secret_value = "val2"
    is_secret    = true
  }
}

resource "azuredevops_check_branch_control" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_variable_group.example.id
  target_resource_type = "variablegroup"
  allowed_branches     = "refs/heads/main, refs/heads/features/*"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID.
* `target_resource_id` - (Required) The ID of the resource being protected by the check.
* `target_resource_type` - (Required) The type of resource being protected by the check. Valid values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`.
* `display_name` - (Required) The name of the branch control check displayed in the web UI.
* `allowed_branches` - (Optional) The branches allowed to use the resource. Specify a comma separated list of allowed branches in `refs/heads/branch_name` format. To allow deployments from all branches, specify ` * ` . `refs/heads/features/* , refs/heads/releases/*` restricts deployments to all branches under features/ or releases/ . Defaults to `*`.
* `verify_branch_protection` - (Optional) Validate the branches being deployed are protected. Defaults to `false`.
* `ignore_unknown_protection_status` - (Optional) Allow deployment from branches for which protection status could not be obtained. Only relevant when verify_branch_protection is `true`. Defaults to `false`.

---

* `timeout` - (Optional) The timeout in minutes for the branch control check. Defaults to `1440`.

## Attributes Reference

In addition to all arguments above the following attributes are exported:

* `id` - The ID of the check.
* `version` - The version of the check.

## Relevant Links

- [Define approvals and checks](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/approvals?view=azure-devops&tabs=check-pass)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Branch Control Check.
* `read` - (Defaults to 1 minute) Used when retrieving the Branch Control Check.
* `update` - (Defaults to 2 minutes) Used when updating the Branch Control Check.
* `delete` - (Defaults to 2 minutes) Used when deleting the Branch Control Check.

## Import

Importing this resource is not supported.
