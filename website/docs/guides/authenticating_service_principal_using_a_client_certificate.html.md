---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with a Client Certificate"
description: |-
  This guide will cover how to use a certificate to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with a Client Certificate

The Azure DevOps provider supports service principals through a variety of authentication methods, including client certificates.

## Service Principal Configuration

1. Create a Service Principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell) and generate a certificate for it. You do not need to assign the service principal any roles in Azure Ad.

2. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration. They may be provided via the `ARM_TENANT_ID` and `ARM_CLIENT_ID` environment variables, or in the provider configuration block with the `tenant_id` and `client_id` attributes.

The certificate may be provided as a base64 string, or by a file on the filesystem with the `ARM_CLIENT_CERTIFICATE` or `ARM_CLIENT_CERTIFICATE_PATH` environment variables, or in the provider configuration block with the `client_certificate` or `client_certificate_path` attributes. To use powershell to base64 encode a .pfx file use `[convert]::ToBase64String((Get-Content -path "cert_with_private_key.pfx" -Encoding byte))`. Note that base64 is **NOT** a security function, and the base64 string should be handled with the same precautions as the original file.

A certificate password may be specified with the `ARM_CLIENT_CERTIFICATE_PASSWORD` environment variable, or in the provider configuration block with the `client_certificate_password` attribute.

### Providing the certificate through the file system

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
  org_service_url = "https://dev.azure.com/my-org"
  client_id                   = "00000000-0000-0000-0000-000000000001"
  tenant_id                   = "00000000-0000-0000-0000-000000000001"
  client_certificate_path     = "C:\\cert.pfx"
  client_certificate_password = "cert password"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### Providing the certificate as a base64 encoded string

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
  org_service_url = "https://dev.azure.com/my-org"
  client_id                   = "00000000-0000-0000-0000-000000000001"
  tenant_id                   = "00000000-0000-0000-0000-000000000001"
  client_certificate          = "MII....lots.and.lots.of.ascii.characters"
  client_certificate_password = "cert password"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```
