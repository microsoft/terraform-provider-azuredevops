---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_required_template"
description: |-
  Manages a Required Template Check.
---

# azuredevops_check_required_template

Manages a Required Template Check.

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

resource "azuredevops_check_required_template" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_serviceendpoint_generic.example.id
  target_resource_type = "endpoint"

  required_template {
    repository_type = "azuregit"
    repository_name = "project/repository"
    repository_ref  = "refs/heads/main"
    template_path   = "template/path.yml"
  }
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

resource "azuredevops_check_required_template" "example" {
  project_id           = azuredevops_project.example.id
  target_resource_id   = azuredevops_environment.example.id
  target_resource_type = "environment"

  required_template {
    repository_name = "project/repository"
    repository_ref  = "refs/heads/main"
    template_path   = "template/path.yml"
  }

  required_template {
    repository_name = "project/repository"
    repository_ref  = "refs/heads/main"
    template_path   = "template/alternate-path.yml"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The project ID. Changing this forces a new Required Template Check to be created.

* `target_resource_id` - (Required) The ID of the resource being protected by the check. Changing this forces a new Required Template Check to be created.

* `target_resource_type` - (Required) The type of resource being protected by the check. Valid values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`. Changing this forces a new Required Template Check to be created.

* `required_template` - (Required) One or more `required_template` blocks documented below.

---

A `required_template` block supports the following:

* `template_path` - (Required) The path to the template yaml.

* `repository_name` - (Required) The name of the repository storing the template.

* `repository_ref` - (Required) The branch in which the template will be referenced.

* `repository_type` - (Optional) The type of the repository storing the template. Possible values are: `azuregit`, `github`, `githubenterprise`, `bitbucket`. Defaults to `azuregit`.

## Attributes Reference

In addition to the Arguments listed above - the following attribute are exported:

* `id` - The ID of the check.
* `version` - The version of the check.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeout) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Required Template Check.
* `read` - (Defaults to 1 minute) Used when retrieving the Required Template Check.
* `update` - (Defaults to 2 minutes) Used when updating the Required Template Check.
* `delete` - (Defaults to 2 minutes) Used when deleting the Required Template Check.

## Import

Importing this resource is not supported.
