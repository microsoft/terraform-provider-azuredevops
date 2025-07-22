---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with a Client Secret"
description: |-
  This guide will cover how to use a client secret to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with a Client Secret

The Azure DevOps provider supports service principals through a variety of authentication methods, including client secrets.

## Service Principal Configuration

1. Create a service principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell). Ignore steps about application roles and certificates.

2. [Generate a client secret for the service principal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-2-create-a-new-application-secret)

3. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration. They may be provided via the `ARM_TENANT_ID` and `ARM_CLIENT_ID` environment variables, or in the provider configuration block with the `tenant_id` and `client_id` attributes.

The client secret may be provided as a string, or by a file on the filesystem with the `ARM_CLIENT_SECRET` or `ARM_CLIENT_SECRET_PATH` environment variables, or in the provider configuration block with the `client_secret` or `client_secret_path` attributes.

### Providing the secret through the file system

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url = "https://dev.azure.com/my-org"
  client_id          = "00000000-0000-0000-0000-000000000001"
  tenant_id          = "00000000-0000-0000-0000-000000000001"
  client_secret_path = "C:\\my_secret.txt"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### Providing the secret directly as a string

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url = "https://dev.azure.com/my-org"
  client_id     = "00000000-0000-0000-0000-000000000001"
  tenant_id     = "00000000-0000-0000-0000-000000000001"
  client_secret = "top-secret-password-string"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```
