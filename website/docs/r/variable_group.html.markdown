---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_variable_group"
description: |-
  Manages variable groups within Azure DevOps project.
---

# azuredevops_variable_group

Manages variable groups within Azure DevOps.

~> **Note**
If Variable Group is linked to a Key Vault, only top 500 secrets will be read by default. Key Vault does not support filter the secret by name, 
we can only read the secrets and do filter in Terraform.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
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
```

## Example Usage With AzureRM Key Vault

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example AzureRM"
  description           = "Managed by Terraform"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Example Subscription Name"
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Variable Group"
  description  = "Example Variable Group Description"
  allow_access = true

  key_vault {
    name                = "example-kv"
    service_endpoint_id = azuredevops_serviceendpoint_azurerm.example.id
  }

  variable {
    name = "key1"
  }

  variable {
    name = "key2"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `name` - (Required) The name of the Variable Group.
- `description` - (Optional) The description of the Variable Group.
- `allow_access` - (Required) Boolean that indicate if this variable group is shared by all pipelines of this project.
- `variable` - (Required) One or more `variable` blocks as documented below.
- `key_vault` -(Optional) A list of `key_vault` blocks as documented below.

A `variable` block supports the following:

- `name` - (Required) The key value used for the variable. Must be unique within the Variable Group.
- `value` - (Optional) The value of the variable. If omitted, it will default to empty string.
- `secret_value` - (Optional) The secret value of the variable. If omitted, it will default to empty string. Used when `is_secret` set to `true`.
- `is_secret` - (Optional) A boolean flag describing if the variable value is sensitive. Defaults to `false`.

A `key_vault` block supports the following:

- `name` - The name of the Azure key vault to link secrets from as variables.
- `service_endpoint_id` - The id of the Azure subscription endpoint to access the key vault.
- `search_depth` - Set the Azure Key Vault Secret search depth. Defaults to `20`. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Variable Group returned after creation in Azure DevOps.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-7.0)
- [Azure DevOps Service REST API 7.0 - Authorized Resources](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/authorizedresources?view=azure-devops-rest-7.0)

## Import
**Variable groups containing secret values cannot be imported.**

Azure DevOps Variable groups can be imported using the project name/variable group ID or by the project Guid/variable group ID, e.g.

```sh
terraform import azuredevops_variable_group.example "Example Project/10"
```

or

```sh
terraform import azuredevops_variable_group.example 00000000-0000-0000-0000-000000000000/0
```

_Note that for secret variables, the import command retrieve blank value in the tfstate._

## PAT Permissions Required

- **Variable Groups**: Read, Create, & Manage
- **Build**: Read & execute
- **Project and Team**: Read
- **Token Administration**: Read & manage
- **Tokens**: Read & manage
- **Work Items**: Read
