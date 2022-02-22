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
resource "azuredevops_project" "test" {
  name = "Test Project"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "Test Variable Group"
  description  = "Test Variable Group Description"
  allow_access = true

  variable {
    name  = "key"
    value = "value"
  }

  variable {
    name         = "Account Password"
    secret_value = "p@ssword123"
    is_secret    = true
  }
}
```

## Example Usage With AzureRM Key Vault

```hcl
resource "azuredevops_project" "test" {
  name = "Test Project"
}

resource "azuredevops_serviceendpoint_azurerm" "test" {
  project_id                = azuredevops_project.test.id
  service_endpoint_name     = "Sample AzureRM"
  description               = "Managed by Terraform"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Sample Subscription"
}

resource "azuredevops_variable_group" "variablegroup" {
  project_id   = azuredevops_project.test.id
  name         = "Test Variable Group"
  description  = "Test Variable Group Description"
  allow_access = true

  key_vault {
    name                = "test-kv"
    service_endpoint_id = azuredevops_serviceendpoint_azurerm.test.id
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

- `project_id` - (Required) The project ID or project name.
- `name` - (Required) The name of the Variable Group.
- `description` - (Optional) The description of the Variable Group.
- `allow_access` - (Required) Boolean that indicate if this variable group is shared by all pipelines of this project.
- `variable` - (Optional) One or more `variable` blocks as documented below.

A `variable` block supports the following:

- `name` - (Required) The key value used for the variable. Must be unique within the Variable Group.
- `value` - (Optional) The value of the variable. If omitted, it will default to empty string.
- `secret_value` - (Optional) The secret value of the variable. If omitted, it will default to empty string. Used when `is_secret` set to `true`.
- `is_secret` - (Optional) A boolean flag describing if the variable value is sensitive. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Variable Group returned after creation in Azure DevOps.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-6.0)
- [Azure DevOps Service REST API 6.0 - Authorized Resources](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/authorizedresources?view=azure-devops-rest-6.0)

## Import
**Variable groups containing secret values cannot be imported.**

Azure DevOps Variable groups can be imported using the project name/variable group ID or by the project Guid/variable group ID, e.g.

```sh
terraform import azuredevops_variable_group.variablegroup "Test Project/10"
```

or

```sh
terraform import azuredevops_variable_group.variablegroup 00000000-0000-0000-0000-000000000000/0
```

_Note that for secret variables, the import command retrieve blank value in the tfstate._

## PAT Permissions Required

~> **Note** After upgrading the API to v6, creating Variable Group linked to Key Vault requires full access permission or you wil get a 401 error.

- **Variable Groups**: Read, Create, & Manage
- **Build**: Read & execute
- **Project and Team**: Read
- **Token Administration**: Read & manage
- **Tokens**: Read & manage
- **Work Items**: Read
