---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating via Managed Identity"
description: |-
  This guide will cover how to use a managed identity to authenticate to Azure DevOps.
---


## What is a managed identity?

[Managed identities for Azure resources](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) can be used to authenticate to services that support Azure Active Directory (Azure AD) authentication. There are two types of managed identities: system-assigned and user-assigned. This article is based on system-assigned managed identities.

Managed identities work in conjunction with Azure Resource Manager (ARM), Azure AD, and the Azure Instance Metadata Service (IMDS). Azure resources that support managed identities expose an internal IMDS endpoint that the client can use to request an access token. No credentials are stored on the VM, and the only additional information needed to bootstrap the Terraform connection to Azure is the subscription ID and tenant ID.

Azure AD creates an AD identity when you configure an Azure resource to use a system-assigned managed identity. The configuration process is described in more detail, below. Azure AD then creates a service principal to represent the resource for role-based access control (RBAC) and access control (IAM). The lifecycle of a system-assigned identity is tied to the resource it is enabled for: it is created when the resource is created and it is automatically removed when the resource is deleted.

Before you can use the managed identity, it has to be configured. [Add the identity to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Configuring Terraform to use a managed identity

Terraform can be configured to use managed identity for authentication in one of two ways: using environment variables, or by defining the fields within the provider block.

### Configuring with environment variables

The `use_msi` must be set to `true` to use a managed identity. By default, Terraform will use the system assigned identity for authentication. To use a user assigned identity instead, you will need to specify the `ARM_CLIENT_ID` environment variable (equivalent to provider block argument `client_id`) to the [client id](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/user_assigned_identity#client_id) of the identity.

A provider block is _technically_ optional when using environment variables. Even so, we recommend defining provider blocks so that you can pin or constrain the version of the provider being used, and configure other optional settings:

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

provider "azuredevops" {
}
```

### Configuring with the provider block

It's also possible to configure a managed identity within the provider block:

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

provider "azuredevops" {
  use_msi = true
}
```
